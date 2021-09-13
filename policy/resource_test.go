package policy

import "testing"

func TestResourceEvent(t *testing.T) {
	//: 6
	resource := MergeRequest
	got := resource.Validate()
	if got != nil {
		t.Errorf("expected nil as `%s` is a valid resource, but received an error %v", resource, got)
	}
}

func TestInvalidResourceEvent(t *testing.T) {
	//: 6
	resource := EventType("invalid")
	got := resource.Validate()
	if got == nil {
		t.Errorf("expected an error as `%s` is an invalid resource, but received no error %v", resource, got)
	}
}
