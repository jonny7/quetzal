//go:build ignore

// this go.gen file generates the requirements-matrix required by GitLab

package main

import (
	"log"
	"os"
	"text/template"
)

type Requirement struct {
	ID       string
	TestCase string
}

func main() {
	f, err := os.Create("../../requirements-matrix.csv")
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("requirements-matrix.csv could not be closed")
		}
	}(f)

	requirements := []Requirement{
		{
			ID:       "REQ-1",
			TestCase: "TestSomething",
		},
		{
			ID:       "REQ-2",
			TestCase: "TestSomethingElse",
		},
	}

	err = packageTemplate.Execute(f, requirements)
	if err != nil {
		log.Println("template failed")
	}
}

var packageTemplate = template.Must(template.New("").Parse(
	`Requirement ID, Test Case{{range .}}
{{.ID}},{{.TestCase}}{{end}}`))
