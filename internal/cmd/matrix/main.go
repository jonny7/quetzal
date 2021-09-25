/*
This is purely a utility function for working with GitLab requirements.

The `matrix` flag will find all the files that
contain `_test.go` as per Go convention and parse them, extracting the function name only (no signature details).
It will look for the special `//: n,n...` where you should specify what requirement(s) in GitLab this test satisfies.
multiple references should be use comma separators with no space.

The `results` flag will generate the requirements.json file used by GitLab to determine if a requirement has been satisfied.
It does this by cross-referencing all the tests for a given requirement, against the tests associated with that requirement from
the matrix.
*/
package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type Requirement struct {
	TestCase string
	ID       int
}

const (
	passed = "passed"
	failed = "failed"
)

type Status struct {
	RequirementID int
	Status        string
}

type JunitResult struct {
	Name    string `xml:"name,attr"`
	Failure string `xml:"failure"`
}

var resultTemplate = "{{$lastIdx := LastIdx (len .)}}{\n\t{{range $i, $el := .}}  \"{{$el.RequirementID}}\":\"{{$el.Status}}\"{{if eq $lastIdx $i}}{{else}},\n{{end}}{{end}}\n}"
var lister = template.Must(template.New("foo").Funcs(template.FuncMap{
	"LastIdx": func(size int) int { return size - 1 },
}).Parse(resultTemplate))

var matrixTemplate = template.Must(template.New("").Parse(
	`Requirement ID, Test Case{{range .}}
REQ-{{.ID}},{{.TestCase}}{{end}}`))

func main() {
	matrix := flag.Bool("matrix", true, "generate the requirements matrix - see `generateMatrix()`")
	results := flag.Bool("results", true, "generate the corresponding results.json file required by GitLab to mark requirements as satisfied or not")

	var reqs []Requirement
	if *matrix {
		reqs = generateMatrix()
	}
	if *results {
		generateResults(reqs)
	}
}

func generateMatrix() []Requirement {
	f, err := os.Create("./requirements-matrix.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("requirements-matrix.csv could not be closed")
		}
	}(f)

	var files []string
	err = filepath.Walk("../../", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "_test.go") {
			files = append(files, path)
		}
		return nil
	})

	var requirements []Requirement
	for _, file := range files {
		contents, fe := os.ReadFile(file)
		if fe != nil {
			log.Fatal(fe)
		}
		reqs := parseFile(string(contents))
		for _, req := range reqs {
			requirements = append(requirements, req)
		}
	}

	sort.Slice(requirements, func(i, j int) bool { return requirements[i].ID < requirements[j].ID })

	err = matrixTemplate.Execute(f, requirements)
	if err != nil {
		log.Println("template failed")
	}
	return requirements
}

func parseFile(contents string) []Requirement {
	re := regexp.MustCompile(`(?:func)\s(Test\w+)(:?.+)(?s:.{0,2}\s+\/\/\:\s?)([0-9]+(,[0-9]+)*)`)
	tests := re.FindAllStringSubmatch(contents, -1)
	var results []Requirement
	for _, test := range tests {
		testName := test[1]
		idSlice := strings.Split(test[3], ",")
		var ids []int
		for _, val := range idSlice {
			valToInt, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, valToInt)
		}
		for _, reqId := range ids {
			results = append(results, Requirement{
				TestCase: testName,
				ID:       reqId,
			})
		}
	}
	return results
}

func generateResults(reqs []Requirement) {
	f, err := os.Create("./requirements.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("requirements.json could not be closed")
		}
	}(f)

	var mappedRequirements = make(map[int][]string)
	for _, r := range reqs {
		mappedRequirements[r.ID] = append(mappedRequirements[r.ID], r.TestCase)
	}

	junit, xmlErr := os.Open("./report.xml")
	if xmlErr != nil {
		log.Fatal(xmlErr)
	}

	var junitResults []JunitResult

	stream := xml.NewDecoder(junit)
	for {
		var junitResult JunitResult
		token, tokenErr := stream.Token()
		if tokenErr != nil {
			break
		}
		switch elementType := token.(type) {
		case xml.StartElement:
			if elementType.Name.Local == "testcase" {
				if e := stream.DecodeElement(&junitResult, &elementType); e != nil {
					log.Fatal(e)
				}
				junitResults = append(junitResults, junitResult)
			}
		}
	}
	var requirementResults []Status
	for k, v := range mappedRequirements {
		status := Status{RequirementID: k}
		var pass = false
		for _, testCase := range v {
			r := filterJunitResult(testCase, junitResults)
			if r == false {
				pass = false
				break
			} else {
				pass = true
			}
		}
		if pass {
			status.Status = passed
		} else {
			status.Status = failed
		}
		requirementResults = append(requirementResults, status)
	}
	err = lister.Execute(f, requirementResults)
	if err != nil {
		log.Println("result template failed")
	}
}

func filterJunitResult(testName string, results []JunitResult) bool {
	for _, j := range results {
		if testName == j.Name {
			if j.Failure == "" {
				return true
			}
		}
	}
	return false
}
