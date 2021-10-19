package main

// Version houses the current version for quetzal
type Version struct {
	version string
}

var current = Version{version: "0.6.0"}

func (v *Version) toString() string {
	return v.version
}
