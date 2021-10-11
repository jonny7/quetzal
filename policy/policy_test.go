package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

// Policy conforms to the matcher interface and is parent
// function that runs each condition and sub-conditions matcher function
// When a new matching element is added, say `Labels`.
// This will need to be updated to reflect that
func TestPolicyMatcher(t *testing.T) {
	//: 3
	policy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}}
	data := []struct {
		name     string
		hook     Webhook
		policy   Policy
		expected bool
		errMsg   string
	}{
		{name: "Matched Policy", hook: Webhook{EventType: gitlab.EventTypeMergeRequest}, policy: policy, expected: true, errMsg: "expected true as policy and hook have the same resource type"},
		{name: "Unmatched Policy", hook: Webhook{EventType: gitlab.EventTypeIssue}, policy: policy, expected: false, errMsg: "expected false as policy and hook do not have the same resource type"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.policy.matcher(d.hook)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}
