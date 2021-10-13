package main

// Version houses the current version for quetzal
type Version struct {
	version string
}

var current = Version{version: "0.3.1"}

func (v *Version) toString() string {
	return v.version
}
