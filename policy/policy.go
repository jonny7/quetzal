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
	Name       string    `yaml:"name,omitempty"`
	Resource   Resource  `yaml:",inline"`
	Conditions Condition `yaml:"conditions,omitempty"`
	Limit      *Limit    `yaml:"limit,omitempty"`
	Actions    *Action   `yaml:"actions,omitempty"`
}

// Resource embeds a gitlab.EventType
type Resource struct {
	EventType gitlab.EventType `yaml:"resource"`
}

// validate confirms that a user specified resource is a valid type
func (r Resource) validate() error {
	switch r.EventType {
	case gitlab.EventTypeBuild, gitlab.EventTypeDeployment, gitlab.EventTypeIssue, gitlab.EventConfidentialIssue, gitlab.EventTypeJob, gitlab.EventTypeMergeRequest, gitlab.EventTypeNote, gitlab.EventConfidentialNote, gitlab.EventTypePipeline, gitlab.EventTypePush, gitlab.EventTypeRelease, gitlab.EventTypeSystemHook, gitlab.EventTypeTagPush, gitlab.EventTypeWikiPage:
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

// conditionMet returns whether this webhook matches a policy resource type
func (r Resource) conditionMet(event GitLabAdaptor) bool {
	return r.EventType == event.ResourceType()
}

// Validate houses a series of validation checks on the
// user specified yaml.
func (p *Policy) Validate() <-chan error {
	ch := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.Resource.validate(); err != nil {
			ch <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.Conditions.Date.validateAll(); err != nil {
			ch <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.Conditions.State.validate(p.Resource.EventType); err != nil {
			ch <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.Conditions.Milestone.validate(); err != nil {
			ch <- err
		}
	}()
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

// ConditionsMet runs a series of checks against all the other conditions that make up a Policy
// in order to report back whether a Policy's criteria is matched by the webhook and an action should occur
func (p Policy) ConditionsMet(event GitLabAdaptor) <-chan Policy {
	valid := make(chan Policy)
	checked := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !p.Resource.conditionMet(event) {
			checked <- false
			return
		}
		checked <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !p.Conditions.State.conditionMet(event) {
			checked <- false
			return
		}
		checked <- true
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
