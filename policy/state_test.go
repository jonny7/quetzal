package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestStateFieldValidation(t *testing.T) {
	//: 11
	data := []struct {
		name          string
		state         *State
		eventType     gitlab.EventType
		expectedIsNil bool
		errMsg        string
	}{
		{name: "No State Listed", state: nil, eventType: gitlab.EventTypeMergeRequest, expectedIsNil: true, errMsg: "expected nil as no state in policy is valid"},
		{name: "Valid MergeEvent State", state: &State{State: []string{string(mergeRequestStateOpen)}}, eventType: gitlab.EventTypeMergeRequest, expectedIsNil: true, errMsg: "expected nil as a MergeEvent can have a state of string(mergeRequestStateOpen)"},
		{name: "Invalid MergeEvent State", state: &State{State: []string{"invalid"}}, eventType: gitlab.EventTypeMergeRequest, expectedIsNil: false, errMsg: "expected an error as a MergeEvent cannot have a state of invalid"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.state.validate(d.eventType)
			want := expectedIsNil(got)
			if want != d.expectedIsNil {
				t.Errorf(d.errMsg)
			}
		})
	}
}
