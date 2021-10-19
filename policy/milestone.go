package policy

// Milestone represents the id of the milestone from GitLab
type Milestone struct {
	Milestone int `yaml:"milestone"`
}

func (m *Milestone) milestone() int {
	if m == nil {
		return 0
	}
	return m.Milestone
}
