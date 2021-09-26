package policy

type Milestone struct {
	Milestone int `yaml:"milestone"`
}

func (m *Milestone) conditionMet(event GitLabAdaptor) bool {
	if m == nil {
		return true
	}
	currentMilestone, err := event.Milestone()
	if err != nil {
		return false // state is being applied to an event that doesn't have a state
	}
	if *currentMilestone == 0 {
		return true
	}
	return *currentMilestone == m.Milestone
}
