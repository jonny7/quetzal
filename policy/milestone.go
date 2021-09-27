package policy

// Milestone represents the integer id from GitLab
type Milestone struct {
	Milestone int `yaml:"milestone"`
}

func (m *Milestone) conditionMet(event GitLabAdaptor) bool {
	if m == nil {
		return true
	}
	currentMilestone, err := event.Milestone()
	if err != nil {
		return false
	}
	if *currentMilestone == 0 {
		return true
	}
	return *currentMilestone == m.Milestone
}
