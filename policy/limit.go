package policy

// Limit is the amount of results to return
type Limit struct {
	MostRecent *int `yaml:"mostRecent,omitempty"`
	Oldest     *int `yaml:"oldest,omitempty"`
}
