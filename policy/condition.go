package policy

// Condition declares what properties and states are required by
// the webhook to have an action performed on it
type Condition struct {
	// Date is a struct to manage date related entries
	Date *Date `yaml:"date,omitempty"`
	// State is the expected state of the webhook event
	State *State `yaml:",inline,omitempty"`
	// Milestone is the milestone of the issue
	Milestone *Milestone `yaml:",inline,omitempty"`
	// Labels provides an array of required labels for the condition to be met
	Labels Labels `yaml:",inline,omitempty"`
	// ForbiddenLabels is an array of labels to not trigger the condition
	ForbiddenLabels ForbiddenLabels `yaml:",inline,omitempty"`
	// Discussion provides a struct to manage whether certain discussion properties meet the given condition
	Discussion *Discussion `yaml:"discussion,omitempty"`
	// Note is the contents of a given note/comment on various different events like commit, mr, issue, code snippet
	Note *Note `yaml:"note"`
}
