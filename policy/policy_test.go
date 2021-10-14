package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

// similar to above this should be updated when more functionality is added to give a more extensive check of
// policy validation
func TestPolicyValidator(t *testing.T) {
	data := []struct {
		name          string
		policy        Policy
		expectedIsNil bool
		errMsg        string
	}{
		{name: "Valid Policy Resource Only", policy: Policy{Resource: Resource{gitlab.EventTypeMergeRequest}}, expectedIsNil: true, errMsg: "expected nil as policy and hook have the same resource type"},
		{name: "Invalid Policy Resource Only", policy: Policy{Resource: Resource{"invalidEntry"}}, expectedIsNil: false, errMsg: "expected an error to be returned as invalid resource type used"},
		{name: "Valid Policy State Only", policy: Policy{Conditions: Condition{State: &State{string(mergeRequestStateApproved)}}}, expectedIsNil: false, errMsg: "expected an error to be returned as invalid resource type used"},
		{name: "Invalid Policy State Only", policy: Policy{Conditions: Condition{State: &State{"invalid"}}}, expectedIsNil: false, errMsg: "expected an error to be returned as invalid resource type used"},
		{name: "Invalid State On Type", policy: Policy{Resource: Resource{gitlab.EventTypeWikiPage}, Conditions: Condition{State: &State{"invalid"}}}, expectedIsNil: false, errMsg: "expected an error to be returned as Wiki Events do not have a state"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.policy.Validate()
			want := expectedIsNil(got)
			if want != d.expectedIsNil {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestPolicyMatcherWithState(t *testing.T) {
	//: 12
	state := &State{string(mergeRequestStateApproved)}
	resource := Resource{gitlab.EventTypeMergeRequest}
	mergeEvent := gitlab.MergeEvent{}
	mergeEvent.ObjectAttributes.Action = string(mergeRequestStateApproved)

	mergeEventUnMatched := gitlab.MergeEvent{}
	mergeEventUnMatched.ObjectAttributes.Action = string(mergeRequestStateClose)

	data := []struct {
		name      string
		policy    Policy
		eventType gitlab.EventType
		adaptor   GitLabAdaptor
		expected  bool
		errMsg    string
	}{
		{name: "Matched Policy Resource", policy: Policy{Resource: resource}, eventType: gitlab.EventTypeMergeRequest, adaptor: MergeEventAdaptor{mergeEvent}, expected: true, errMsg: "expected true as the policy and hook match on resource"},
		{name: "Unmatched Policy Resource", policy: Policy{Resource: resource}, eventType: gitlab.EventTypeWikiPage, adaptor: MergeEventAdaptor{mergeEvent}, expected: false, errMsg: "expected false as the policy and hook do not match on resource"},
		{name: "Matched Policy Resource and Nil State", policy: Policy{Resource: resource}, eventType: gitlab.EventTypeMergeRequest, adaptor: MergeEventAdaptor{mergeEvent}, expected: true, errMsg: "expected true as the policy and hook match on state"},
		{name: "Matched Policy Resource and State", policy: Policy{Resource: resource, Conditions: Condition{State: state}}, eventType: gitlab.EventTypeMergeRequest, adaptor: MergeEventAdaptor{mergeEvent}, expected: true, errMsg: "expected true as the policy and hook match on state"},
		{name: "Unmatched Policy Resource and State", policy: Policy{Resource: resource, Conditions: Condition{State: state}}, eventType: gitlab.EventTypeMergeRequest, adaptor: MergeEventAdaptor{mergeEventUnMatched}, expected: false, errMsg: "expected false as the policy and hook do not match on state"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := matcher(d.policy, d.adaptor, d.eventType)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}
}

func TestLabelMatching(t *testing.T) {
	//: 15
	adaptor := MergeEventAdaptor{}
	adaptor.Labels = append(adaptor.Labels, &gitlab.Label{Name: "needed"})
	adaptor2Labels := adaptor
	adaptor2Labels.Labels = append(adaptor2Labels.Labels, &gitlab.Label{Name: "important"})
	policy := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}}
	policyWithLabel := policy
	policyWithLabel.Conditions.Labels = []string{"needed"}
	policy3 := Policy{Resource: Resource{gitlab.EventTypeMergeRequest}, Conditions: Condition{Labels: []string{"important", "needed"}}}

	data := []struct {
		name      string
		policy    Policy
		eventType gitlab.EventType
		adaptor   GitLabAdaptor
		expected  bool
		errMsg    string
	}{
		{name: "Policy with No Labels", policy: policy, eventType: gitlab.EventTypeMergeRequest, adaptor: adaptor, expected: true, errMsg: "expected true as the policy has no required labels"},
		{name: "Policy and Webhook with same 1 Label", policy: policyWithLabel, eventType: gitlab.EventTypeMergeRequest, adaptor: adaptor, expected: true, errMsg: "expected true as the policy and hook have the same labels"},
		{name: "Policy with 1 label, webhook with 2", policy: policyWithLabel, eventType: gitlab.EventTypeMergeRequest, adaptor: adaptor2Labels, expected: false, errMsg: "expected false as the policy and hook only partially match"},
		{name: "Policy with 2 label, webhook with 2, different orders", policy: policy3, eventType: gitlab.EventTypeMergeRequest, adaptor: adaptor2Labels, expected: true, errMsg: "expected true as the policy and hook have the same slices, just in different orders"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := matcher(d.policy, d.adaptor, d.eventType)
			if got != d.expected {
				t.Errorf(d.errMsg)
			}
		})
	}

}
