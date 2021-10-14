package policy

import (
	"github.com/xanzy/go-gitlab"
	"strings"
)

// fieldValidator ensures that a field or entire struct
// has valid user specified input
type fieldValidator interface {
	validate() error
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

func (p Policy) resource() gitlab.EventType {
	return p.Resource.EventType
}

func (p Policy) state() *string {
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
	return sliceLower(p.Conditions.Labels)
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
	State *State `yaml:",inline,omitempty"`
	// Milestone is the milestone of the issue
	Milestone *Milestone `yaml:",inline,omitempty"`
	// Labels provides an array of required labels for the condition to be met
	Labels []string `yaml:"labels"`
	// ForbiddenLabels is an array of labels to not trigger the condition
	ForbiddenLabels []string `yaml:"forbiddenLabels"`
	// Discussion provides a struct to manage whether certain discussion properties meet the given condition
	//Discussion *Discussion `yaml:"discussion,omitempty"` @todo
	// Note is the contents of a given note/comment on various different events like commit, mr, issue, code snippet
	Note *Note `yaml:"note"`
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

const (
	createdAt DateAttribute = "created_at"
	updatedAt DateAttribute = "updated_at"
)

// DateCondition is the greater than or less than [date] filter
type DateCondition string

const (
	olderThan DateCondition = "older_than"
	newerThan DateCondition = "newer_than"
)

// DateIntervalType is the type of available interval
type DateIntervalType string

const (
	days   DateIntervalType = "days"
	weeks  DateIntervalType = "weeks"
	months DateIntervalType = "months"
	years  DateIntervalType = "years"
)

// issueState represents the possible states an issue can be in
type issueState string

const (
	issueStateOpen   issueState = "open"
	issueStateClose  issueState = "close"
	issueStateReopen issueState = "reopen"
	issueStateUpdate issueState = "update"
)

// releaseState represents the possible states an releaseState can be in
type releaseState string

const (
	releaseStateCreate releaseState = "create"
	releaseStateUpdate releaseState = "update"
)

// Milestone represents the integer id from GitLab
type Milestone struct {
	Milestone int `yaml:"milestone"`
}

// NoteType is the type of note: Commit, MergeRequest, Issue, Snippet
type NoteType string

// Mentions is an array of users mentioned in a comment
type Mentions []string

// Command is a string backed type for a given command to respond to
type Command string

const (
	// NoteCommit are comments on Commits
	NoteCommit NoteType = "Commit"
	// NoteMergeRequest are comments on MergeRequests
	NoteMergeRequest NoteType = "MergeRequest"
	// NoteIssue are comments on Issues
	NoteIssue NoteType = "Issue"
	// NoteSnippet are comments on Snippets
	NoteSnippet NoteType = "Snippet"
)

// Note represents a GitLab Note, which is essentially a comment on
// a series of different scenarios and event types
type Note struct {
	// Type is the NoteType of the note from GitLab. If you need to narrow down
	// the type of note then use this, if left blank, then it will apply to all note types
	Type *NoteType `yaml:"noteType"`
	// Mentions looks for user's mentioned in the note
	Mentions Mentions `yaml:"mentions"`
	// Command is the specified string to look for if needed.
	Command Command `yaml:"command"`
}
