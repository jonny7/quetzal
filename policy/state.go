package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

// State represents the webhook state, this is only available
// on certain events
type State struct {
	State string `yaml:"state"`
}

// mergeRequestState represents the possible states a merge request can be in
type mergeRequestState string

const (
	mergeRequestStateOpen       mergeRequestState = "open"
	mergeRequestStateClose      mergeRequestState = "close"
	mergeRequestStateReopen     mergeRequestState = "reopen"
	mergeRequestStateUpdate     mergeRequestState = "update"
	mergeRequestStateApproved   mergeRequestState = "approved"
	mergeRequestStateUnApproved mergeRequestState = "unapproved"
	mergeRequestStateMerge      mergeRequestState = "merge"
)

func (s *State) state() *string {
	if s == nil {
		return nil
	}
	return &s.State
}

func (s *State) validate(eventType gitlab.EventType) error {
	if s == nil {
		return nil
	}
	if eventType == gitlab.EventTypeMergeRequest {
		return validateMergeRequestState(*s)
	}
	return fmt.Errorf("the state condition was used on an unexpected event type :%s", eventType)
}

// validates that a given state for MergeEvents is valid
func validateMergeRequestState(s State) error {
	switch mergeRequestState(strings.ToLower(s.State)) {
	case mergeRequestStateOpen, mergeRequestStateClose, mergeRequestStateReopen, mergeRequestStateUpdate, mergeRequestStateApproved, mergeRequestStateUnApproved, mergeRequestStateMerge:
		return nil
	}
	return fmt.Errorf("available states for Merge Requests are `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`", mergeRequestStateOpen, mergeRequestStateClose, mergeRequestStateReopen, mergeRequestStateUpdate, mergeRequestStateApproved, mergeRequestStateUnApproved, mergeRequestStateMerge)
}
