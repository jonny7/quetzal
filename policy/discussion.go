package policy

// DiscussionAttribute is the type of discussion
type DiscussionAttribute string

// DiscussionCondition identifies to amount of threads required
type DiscussionCondition string

// Discussion represents a Gitlab discussion
type Discussion struct {
	// Attribute can be `threads` or `notes`
	Attribute DiscussionAttribute `yaml:"attribute"`
	// Condition can be `less_than` or `greater_than`
	Condition DiscussionCondition `yaml:"condition"`
	// Threshold is an integer value of how many discussion items
	Threshold int `yaml:"threshold"`
}
