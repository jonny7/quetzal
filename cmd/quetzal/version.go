package main

// Version houses the current version for quetzal
type Version struct {
	version string
}

var current = ""

func getVersion() string {
	v := Version{version: current}
	return v.version
}
