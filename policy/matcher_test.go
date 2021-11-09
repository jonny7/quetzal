package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestMatcherMilestones(t *testing.T) {
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

func TestMatcherForbiddenLabels(t *testing.T) {
	baseAdaptor := MergeEventAdaptor{gitlab.MergeEvent{}}
	forbiddenLabelsPolicy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}, Conditions: Condition{ForbiddenLabels: ForbiddenLabels{[]string{"v1", "v2"}}}}
	forbiddenLabelsAdaptor := baseAdaptor
	forbiddenLabelsAdaptor.Labels = append(forbiddenLabelsAdaptor.Labels, &gitlab.Label{Name: "something"})
	forbiddenLabelsAdaptorPartiallyMissing := baseAdaptor
	forbiddenLabelsAdaptorPartiallyMissing.Labels = append(forbiddenLabelsAdaptor.Labels, &gitlab.Label{Name: "v1"})

	data := []struct {
		name     string
		policy   Matcher
		adaptor  GitLabAdaptor
		expected bool
		errMsg   string
	}{
		{name: "MR missing both forbidden labels", policy: forbiddenLabelsPolicy, adaptor: forbiddenLabelsAdaptor, expected: true, errMsg: "expected true as adaptor doesn't have any of the labels required"},
		{name: "MR missing one forbidden labels", policy: forbiddenLabelsPolicy, adaptor: forbiddenLabelsAdaptorPartiallyMissing, expected: false, errMsg: "expected false as adaptor isn't missing all of the forbidden labels"},
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
