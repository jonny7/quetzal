package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestMatcher(t *testing.T) {
	baseAdaptor := MergeEventAdaptor{gitlab.MergeEvent{}}
	basePolicy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}}
	milestonePolicy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}, Conditions: Condition{Milestone: &Milestone{123}}}
	milestoneAdaptor := baseAdaptor
	milestoneAdaptor.ObjectAttributes.MilestoneID = 123

	data := []struct {
		name     string
		policy   Matcher
		adaptor  GitLabAdaptor
		expected bool
		errMsg   string
	}{
		{name: "No Milestones on adaptor", policy: milestonePolicy, adaptor: baseAdaptor, expected: false, errMsg: "expected false as policy and adaptor do not match"},
		{name: "Milestones match", policy: milestonePolicy, adaptor: milestoneAdaptor, expected: true, errMsg: "expected true as policy and adaptor match"},
		{name: "No Policy milestone", policy: basePolicy, adaptor: baseAdaptor, expected: true, errMsg: "expected true as policy and adaptor don't have milestones"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := matcher(d.policy, d.adaptor, gitlab.EventTypeMergeRequest)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}
