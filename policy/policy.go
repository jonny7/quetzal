package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"sync"
)

// Validator provides a method to validate a struct
// created from a yaml config file, most types
// are backed by strings, which means that value can
// be passed in the .yml file. Validate confirms these are
// permitted values
type Validator interface {
	Validate() error
}

// Policies contains a slice of `Policy`s
type Policies struct {
	Policies []Policy `yaml:"policies"`
}

// Policy is a containing struct that identifies the
// required policy for a certain webhook
type Policy struct {
	Name       string     `yaml:"name,omitempty"`
	Resource   Resource   `yaml:",inline"`
	Conditions *Condition `yaml:"conditions,omitempty"`
	Limit      *Limit     `yaml:"limit,omitempty"`
	Actions    *Action    `yaml:"actions,omitempty"`
}

// Resource embeds a gitlab.EventType
type Resource struct {
	EventType gitlab.EventType `yaml:"resource"`
}

func (r Resource) validate() error {
	switch r.EventType {
	case gitlab.EventTypeBuild,
		gitlab.EventTypeDeployment,
		gitlab.EventTypeIssue,
		gitlab.EventConfidentialIssue,
		gitlab.EventTypeJob,
		gitlab.EventTypeMergeRequest,
		gitlab.EventTypeNote,
		gitlab.EventConfidentialNote,
		gitlab.EventTypePipeline,
		gitlab.EventTypePush,
		gitlab.EventTypeRelease,
		gitlab.EventTypeSystemHook,
		gitlab.EventTypeTagPush,
		gitlab.EventTypeWikiPage:
		return nil
	}
	return fmt.Errorf("`policy:resource` allowed options are: `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`. But received: %v", gitlab.EventTypeBuild,
		gitlab.EventTypeDeployment,
		gitlab.EventTypeIssue,
		gitlab.EventConfidentialIssue,
		gitlab.EventTypeJob,
		gitlab.EventTypeMergeRequest,
		gitlab.EventTypeNote,
		gitlab.EventConfidentialNote,
		gitlab.EventTypePipeline,
		gitlab.EventTypePush,
		gitlab.EventTypeRelease,
		gitlab.EventTypeSystemHook,
		gitlab.EventTypeTagPush,
		gitlab.EventTypeWikiPage, r.EventType)
}

func (r Resource) conditionMet(event GitLabAdaptor) bool {
	return r.EventType == event.ResourceType()
}

// Validate houses a series of validation checks on the
// user specified yaml.
func (p *Policy) Validate() error {
	if err := p.Resource.validate(); err != nil {
		return err
	}
	if p.Conditions.Date != nil {
		if err := p.Conditions.Date.Attribute.validate(); err != nil {
			return err
		}
		if err := p.Conditions.Date.Condition.validate(); err != nil {
			return err
		}
		if err := p.Conditions.Date.IntervalType.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p Policy) ConditionsMet(event GitLabAdaptor) <-chan Policy {
	valid := make(chan Policy)
	checked := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if p.Resource.EventType != event.ResourceType() {
			checked <- false
		}
		checked <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		state, err := event.State()
		if err != nil {
			checked <- false
		}
		checked <- state == p.Conditions.State
	}()

	go func() {
		wg.Wait()
		close(checked)
	}()

	go func(correct bool) {
		defer close(valid)
		for result := range checked {
			if !result {
				correct = false
				break
			}
		}
		if correct {
			valid <- p
		}
	}(true)

	return valid
}
