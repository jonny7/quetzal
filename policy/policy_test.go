package policy

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

func TestConditionMetResourceType(t *testing.T) {
	//: 6
	adaptor := MergeEventAdaptor{gitlab.MergeEvent{
		ObjectKind: string(gitlab.EventTypeMergeRequest)},
	}
	p := Policy{Resource: Resource{
		EventType: gitlab.EventTypeBuild,
	}}

	got := p.Resource.conditionMet(adaptor)
	if got {
		t.Errorf("expected false as resource types don't match.")
	}
}
