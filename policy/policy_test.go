package policy

import (
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
)

func TestDateValidationIntegration(t *testing.T) {
	//: 1
	p := Policy{
		Resource: Resource{
			EventType: gitlab.EventTypeNote,
		},
		Conditions: Condition{
			Date: &Date{
				Attribute: "nothing",
			},
		},
	}

	got := <-p.Validate()
	if got == nil {
		t.Errorf("expected an error for invalid yaml date")
	}
}

func TestResourceValidationIntegration(t *testing.T) {
	//: 6
	p := Policy{
		Resource: Resource{
			EventType: gitlab.EventType("Not A Hook"),
		},
	}

	got := <-p.Validate()
	if !strings.Contains(got.Error(), "`policy:resource` allowed options are:") {
		t.Errorf("expected an error for invalid yaml resource")
	}
}

func TestStateValidationIntegration(t *testing.T) {
	//: 15,17
	p := Policy{
		Resource: Resource{
			EventType: gitlab.EventTypeMergeRequest,
		},
		Conditions: Condition{
			State: &State{"invalid"},
		},
	}

	got := <-p.Validate()
	if !strings.Contains(got.Error(), "available states for Merge Requests are") {
		t.Errorf("expected an error for invalid yaml state")
	}
}

func TestMilestoneValidationIntegration(t *testing.T) {
	//: 18
	p := Policy{
		Resource: Resource{
			EventType: gitlab.EventTypeMergeRequest,
		},
		Conditions: Condition{
			Milestone: &Milestone{0},
		},
	}

	got := <-p.Validate()
	if got == nil {
		t.Errorf("expected an error for invalid yaml state")
	}
}
