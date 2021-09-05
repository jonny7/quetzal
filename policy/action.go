package policy

// Action struct houses how an eligible webhook
// event should be responded to
type Action struct {
	// HTTP contains allows a user to call an API endpoint
	// and use that in subsequent actions
	HTTP HTTP `yaml:"http,omitempty"`
	// Labels identifies which labels to add to an issue
	Labels []string `yaml:"labels,omitempty"`
	// RemoveLabels defines what labels to remove
	RemoveLabels []string `yaml:"remove_labels,omitempty"`
	// Status sets the status of the issue
	Status string `yaml:"status,omitempty"`
	// Mention allows the bot to mention certain users
	Mention []string `yaml:"mention,omitempty"`
	// Comment will leave a comment on said issue
	Comment string `yaml:"comment,omitempty"`
}
