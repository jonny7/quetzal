package policy

import (
	"github.com/xanzy/go-gitlab"
	"strings"
)

// FieldValidator ensures that a field has valid user specified input
type FieldValidator interface {
	fieldValidator(eventType gitlab.EventType) error
}

// Validator allows the Policies to be checked for invalid
// or incompatible instructions
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
	Name       string    `yaml:"name"`
	Resource   Resource  `yaml:",inline"`
	Conditions Condition `yaml:"conditions,omitempty"`
	//Limit      *Limit    `yaml:"limit,omitempty"` @todo
	Actions Action `yaml:"actions,omitempty"`
}

func (p Policy) milestone() int {
	return p.Conditions.Milestone.milestone()
}

func (p Policy) resource() gitlab.EventType {
	return p.Resource.EventType
}

func (p Policy) state() []string {
	return p.Conditions.State.state()
}

// lower a slice so label matching is case-insensitive
func sliceLower(sl []string) []string {
	var lowered []string
	for _, l := range sl {
		lowered = append(lowered, strings.ToLower(l))
	}
	return lowered
}

func (p Policy) labels() []string {
	return sliceLower(p.Conditions.Labels.Labels)
}

func (p Policy) forbiddenLabels() []string {
	return sliceLower(p.Conditions.ForbiddenLabels.ForbiddenLabels)
}

// Validate validates a Policy's correctness
func (p Policy) Validate() error {
	// validate conditions
	if err := p.Resource.validate(); err != nil {
		return err
	}
	if err := p.Conditions.State.validate(p.Resource.EventType); err != nil {
		return err
	}
	// validate actions
	if err := p.Actions.validate(p.Resource.EventType); err != nil {
		return err
	}
	return nil
}

// Condition declares what properties and states are required by
// the webhook to have an action performed on it
type Condition struct {
	// Date is a struct to manage date related entries
	//Date *Date `yaml:"date,omitempty"`
	// State is the expected state of the webhook event
	State *State `yaml:",inline"`
	// Milestone is the milestone of the issue
	Milestone *Milestone `yaml:",inline"`
	// Labels provides an array of required labels for the condition to be met
	Labels Labels `yaml:",inline"`
	// ForbiddenLabels is an array of labels that need to all be missing to
	ForbiddenLabels ForbiddenLabels `yaml:",inline"`
	// Discussion provides a struct to manage whether certain discussion properties meet the given condition
	//Discussion *Discussion `yaml:"discussion,omitempty"` @todo
	// Note is the contents of a given note/comment on various different events like commit, mr, issue, code snippet
	Note *Note `yaml:"note"`
}

// ForbiddenLabels is a list of labels that are missing from an issue and will trigger an action
type ForbiddenLabels struct {
	ForbiddenLabels []string `yaml:"forbiddenLabels"`
}

// Date is possible condition that can be used to allow or
// disallow the behaviour of the Bot see `config.yaml`
type Date struct {
	// Attribute can be `created_at` or `updated_at`
	Attribute DateAttribute `yaml:"attribute"`
	// Condition can be `older_than` or `newer_than`
	Condition DateCondition `yaml:"condition"`
	// IntervalType can be `days`, `weeks`, `months`, `years`
	IntervalType DateIntervalType `yaml:"intervalType"`
	// Interval is a numeric representation of the `IntervalType`
	Interval int `yaml:"interval"`
}

// DateAttribute is the updated or created property
type DateAttribute string

//const (
//	createdAt DateAttribute = "created_at"
//	updatedAt DateAttribute = "updated_at"
//)

// DateCondition is the greater than or less than [date] filter
type DateCondition string

//const (
//	olderThan DateCondition = "older_than"
//	newerThan DateCondition = "newer_than"
//)

// DateIntervalType is the type of available interval
type DateIntervalType string

//const (
//	days   DateIntervalType = "days"
//	weeks  DateIntervalType = "weeks"
//	months DateIntervalType = "months"
//	years  DateIntervalType = "years"
//)

// issueState represents the possible states an issue can be in
//type issueState string

//const (
//	issueStateOpen   issueState = "open"
//	issueStateClose  issueState = "close"
//	issueStateReopen issueState = "reopen"
//	issueStateUpdate issueState = "update"
//)

// releaseState represents the possible states an releaseState can be in
//type releaseState string

//const (
//	releaseStateCreate releaseState = "create"
//	releaseStateUpdate releaseState = "update"
//)
