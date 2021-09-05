package policy

// RequestMethod is the HTTP method type
type RequestMethod string

// HTTP provides a struct to allow yaml generated actions
// for responding to webhooks
type HTTP struct {
	// Once is whether to trigger this action every time the webhook condition is met
	// examples of when you might want to set this to true are you want an initial request
	// to return a series of valid users for approvals and to iterate through them as MRs
	// come in. Each value i
	Once bool `yaml:"once,omitempty"`
	// Method is the type of HTTP request required for the GitLab endpoint
	Method RequestMethod `yaml:"method,omitempty"`
	// Endpoint is the endpoint for GitLab eg /project/:id/approval_merge
	Endpoint string `yaml:"endpoint,omitempty"`
	// Iterate allows a property in the response to be iterated on each subsequent webhook which meets the
	// specified condition.
	Iterate struct {
		Attribute string `yaml:"attribute,omitempty"`
	} `yaml:"iterate,omitempty"`
}
