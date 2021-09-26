package policy

import "strings"

// State represents the webhook state, this is only available
// on certain events
type State struct {
	State string `yaml:"state"`
}

func (s *State) conditionMet(event GitLabAdaptor) bool {
	if s == nil {
		return true // no state on policy, so it can't fail
	}
	currentState, err := event.State()
	if err != nil {
		return false // state is being applied to an event that doesn't have a state
	}
	return strings.ToLower(currentState) == strings.ToLower(s.State)
}
