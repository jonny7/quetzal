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

// issueState represents the possible states an issue can be in
type issueState string

const (
	issueStateOpen   issueState = "open"
	issueStateClose  issueState = "close"
	issueStateReopen issueState = "reopen"
	issueStateUpdate issueState = "update"
)

// releaseState represents the possible states an releaseState can be in
type releaseState string

const (
	releaseStateCreate releaseState = "create"
	releaseStateUpdate releaseState = "update"
)

func (s *State) conditionMet(event GitLabAdaptor) bool {
	if s == nil {
		return true // no state on policy, so it can't fail
	}
	currentState, err := event.State()
	if err != nil {
		return false // state is being applied to an event that doesn't have a state
	}
	return strings.ToLower(*currentState) == strings.ToLower(s.State)
}

func (s *State) validate(p Policy) error {
	if s == nil {
		return nil
	}
	if p.Resource.EventType == gitlab.EventTypeMergeRequest {
		return validateMergeRequestState(*s)
	}
	if p.Resource.EventType == gitlab.EventTypeIssue {
		return validateIssueState(*s)
	}
	if p.Resource.EventType == gitlab.EventTypeRelease {
		return validateReleaseState(*s)
	}
	return fmt.Errorf("state can not be applied on the webhook %s", p.Resource.EventType)
}

func validateMergeRequestState(s State) error {
	switch mergeRequestState(strings.ToLower(s.State)) {
	case mergeRequestStateOpen, mergeRequestStateClose, mergeRequestStateReopen, mergeRequestStateUpdate, mergeRequestStateApproved, mergeRequestStateUnApproved, mergeRequestStateMerge:
		return nil
	}
	return fmt.Errorf("available states for Merge Requests are `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`", mergeRequestStateOpen, mergeRequestStateClose, mergeRequestStateReopen, mergeRequestStateUpdate, mergeRequestStateApproved, mergeRequestStateUnApproved, mergeRequestStateMerge)
}

func validateIssueState(s State) error {
	switch issueState(strings.ToLower(s.State)) {
	case issueStateOpen, issueStateClose, issueStateReopen, issueStateUpdate:
		return nil
	}
	return fmt.Errorf("available states for Issues are `%s`, `%s`, `%s`, `%s`", issueStateOpen, issueStateClose, issueStateReopen, issueStateUpdate)
}

func validateReleaseState(s State) error {
	switch releaseState(strings.ToLower(s.State)) {
	case releaseStateCreate, releaseStateUpdate:
		return nil
	}
	return fmt.Errorf("available states for Releases are `%s`, `%s`", releaseStateCreate, releaseStateUpdate)
}
