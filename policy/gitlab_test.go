package policy

import (
	"testing"
)

func TestEventTypePrepareUpdate(t *testing.T) {
	//: 6
	data := []struct {
		name     string
		event    GitLabAdaptor
		action   Action
		expected int
		errMsg   string
	}{
		{name: "0 prepared fn()", event: MergeEventAdaptor{}, action: Action{}, expected: 0, errMsg: "expected %d funcs to be stacked as action has nothing, but got %d"},
		{name: "1 prepared fn()", event: MergeEventAdaptor{}, action: Action{Labels: Labels{[]string{"added"}}}, expected: 1, errMsg: "expected %d funcs to be stacked as action has content, but got %d"},
		{name: "2 prepared fn()", event: MergeEventAdaptor{}, action: Action{Labels: Labels{[]string{"added"}}, Comment: "This is the func you are looking for"}, expected: 2, errMsg: "expected %d funcs to be stacked as action has content, but got %d"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.event.prepareUpdates(d.action)
			if len(got) != d.expected {
				t.Errorf(d.errMsg, d.expected, len(got))
			}
		})
	}
}
