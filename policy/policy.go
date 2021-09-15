package policy

import "github.com/xanzy/go-gitlab"

// Policy is a containing struct that identifies the
// required policy for a certain webhook
type Policy struct {
	Name       string           `yaml:"name,omitempty"`
	Resource   gitlab.EventType `yaml:"resource"`
	Conditions *Condition       `yaml:"conditions,omitempty"`
	Limit      *Limit           `yaml:"limit,omitempty"`
	Actions    *Action          `yaml:"actions,omitempty"`
}

// Policies contains a slice of `Policy`s
type Policies struct {
	Policies []Policy `yaml:"policies"`
}
