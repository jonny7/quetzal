package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestConditionsMetLabels(t *testing.T) {
	//: 7,19
	adaptor := MergeEventAdaptor{MergeEvent: gitlab.MergeEvent{
		Labels: []*gitlab.Label{{Name: "something"}},
	}}

	p := Policy{
		Conditions: Condition{
			Labels: Labels{labels: []string{"something"}},
		},
	}

	got := p.Conditions.Labels.conditionMet(adaptor)
	if !got {
		t.Errorf("expected true as webhook contained all labels in policy")
	}
}

func TestConditionsMetLabelsNegative(t *testing.T) {
	//: 7,19
	adaptor := MergeEventAdaptor{MergeEvent: gitlab.MergeEvent{
		Labels: []*gitlab.Label{{Name: "something"}},
	}}

	p := Policy{
		Conditions: Condition{
			Labels: Labels{labels: []string{"something", "another"}},
		},
	}

	got := p.Conditions.Labels.conditionMet(adaptor)
	if got != false {
		t.Errorf("expected false as policy has 2 labels and webhook only one")
	}
}

func TestConditionsMetNoLabels(t *testing.T) {
	//: 7,19
	adaptor := MergeEventAdaptor{MergeEvent: gitlab.MergeEvent{
		Labels: []*gitlab.Label{{Name: "something"}},
	}}

	p := Policy{
		Conditions: Condition{},
	}

	got := p.Conditions.Labels.conditionMet(adaptor)
	if got != true {
		t.Errorf("expected true as policy has 0 labels")
	}
}
