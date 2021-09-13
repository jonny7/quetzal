package policy

// Condition declares what properties and states are required by
// the webhook to have an action performed on it
type Condition struct {
	// Date is a struct to manage date related entries
	Date *Date `yaml:"date,omitempty"`
	// State is the state of the webhook issue
	State string `yaml:"state,omitempty"`
	// Milestone is the milestone of the issue
	Milestone string `yaml:"milestone,omitempty"`
	// Labels provides an array of required labels for the condition to be met
	Labels []string `yaml:"labels,omitempty"`
	// ForbiddenLabels is an array of labels to not trigger the condition
	ForbiddenLabels []string `yaml:"forbiddenLabels,omitempty"`
	// Discussion provides a struct to manage whether certain discussion properties meet the given condition
	Discussion *Discussion `yaml:"discussion,omitempty"`
	// Note is the contents of a given note/comment on various different events like commit, mr, issue, code snippet
	Note *Note `yaml:"note"`
}
