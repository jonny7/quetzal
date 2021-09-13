package policy

import (
	"testing"
)

func TestDateIntervalTypeValidate(t *testing.T) {
	//: 3
	dit := DateIntervalType("weeks")
	got := dit.Validate()
	if got != nil {
		t.Errorf("expected nil as `%s` is a valid config entry, but received an error %v", dit, got)
	}
}

func TestDateIntervalTypeValidateInvalid(t *testing.T) {
	//: 3
	dit := DateIntervalType("kittens")
	got := dit.Validate()
	if got == nil {
		t.Errorf("expected an error, but did not recevie one for %s", dit)
	}
}

func TestDateConditionValidate(t *testing.T) {
	//: 2
	dc := DateCondition("older_than")
	got := dc.Validate()
	if got != nil {
		t.Errorf("expected nil as `%s` is a valid entry, but received an error %v", dc, got)
	}
}

func TestDateConditionValidateInvalid(t *testing.T) {
	//: 2
	dc := DateCondition("less_than")
	got := dc.Validate()
	if got == nil {
		t.Errorf("expected an error, but did not recevie one for %s", dc)
	}
}

func TestDateAttributeValidate(t *testing.T) {
	//: 1
	da := DateAttribute("created_at")
	got := da.Validate()
	if got != nil {
		t.Errorf("expected nil as `%s` is a valid entry, but received an error %v", da, got)
	}
}

func TestDateAttributeValidateInvalid(t *testing.T) {
	//: 1
	da := DateAttribute("another_date")
	got := da.Validate()
	if got == nil {
		t.Errorf("expected an error, but did not recevie one for %s", da)
	}
}
