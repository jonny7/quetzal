package policy

import "fmt"

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

func (m *Milestone) validate() error {
	if m == nil || m.Milestone > 0 {
		return nil
	}
	return fmt.Errorf("a milestone of zero, is almost certainly an error. If this wasn't intended, it likely means that your policy provided value was given the zero value of an int. %d", m.Milestone)
}
