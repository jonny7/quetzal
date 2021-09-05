package policy

// Limit is the amount of results to return
type Limit struct {
	MostRecent int `yaml:"most_recent,omitempty"`
	Oldest     int `yaml:"oldest,omitempty"`
}
