package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

// useful if working with tests expecting an error or nil,
// and you don't want to use error strings to check
func expectedIsNil(a interface{}) bool {
	return a == nil
}

// Ensures that user specified policies are valid
func TestResourceFieldValidation(t *testing.T) {
	//: 1
	data := []struct {
		name          string
		resource      Resource
		expectedIsNil bool
		errMsg        string
	}{
		{name: "Invalid Resource", resource: Resource{"not-valid"}, expectedIsNil: false, errMsg: "expected not nil as not-valid is not a valid gitlab.EventType"},
		{name: "Valid Resource", resource: Resource{gitlab.EventTypeMergeRequest}, expectedIsNil: true, errMsg: "expected nil as EventTypeMergeRequest is a valid gitlab.EventType"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.resource.validate()
			want := expectedIsNil(got)
			if want != d.expectedIsNil {
				t.Errorf(d.errMsg)
			}
		})
	}
}

// compares the policy to resource, to the incoming webhook
// and returns whether these match
func TestResourceFieldMatcher(t *testing.T) {
	//: 2
	policy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}}
	data := []struct {
		name     string
		hook     Webhook
		policy   Policy
		expected bool
		errMsg   string
	}{
		{name: "UnMatched Resource", hook: Webhook{EventType: gitlab.EventTypeMergeRequest}, policy: policy, expected: true, errMsg: "expected true as policy hook have the same resource type"},
		{name: "Valid Resource", hook: Webhook{EventType: gitlab.EventTypeIssue}, policy: policy, expected: false, errMsg: "expected false as policy hook do not have the same resource type"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.policy.Resource.matcher(d.hook)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}
