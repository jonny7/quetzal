package policy

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TestDateIntervalTypeValidate(t *testing.T) {
	dit := DateIntervalType("weeks")
	got := dit.Validate()
	if got != nil {
		t.Errorf("expected nil as `%v` is a valid config entry, but received an error %v", dit, got)
	}
}

func TestDateIntervalTypeValidateInvalid(t *testing.T) {
	dit := DateIntervalType("kittens")
	got := dit.Validate()
	if got == nil {
		t.Errorf("expected an error, but received an error %v", got)
	}
}

func TestDateConditionValidate(t *testing.T) {
	dc := DateCondition("older_than")
	got := dc.Validate()
	if got != nil {
		t.Errorf("expected nil as `%v` is a valid entry, but received an error %v", dc, got)
	}
}

func TestDateConditionValidateInvalid(t *testing.T) {
	dc := DateCondition("less_than")
	got := dc.Validate()
	if got == nil {
		t.Errorf("expected an error, but received an error %v", got)
	}
}

func TestDateAttributeValidate(t *testing.T) {
	id := 1
	f := getFunctionName(getFunctionName(TestDateAttributeValidate))
	fmt.Println(id, f)
	da := DateAttribute("created_at")
	got := da.Validate()
	if got != nil {
		t.Errorf("expected nil as `%v` is a valid entry, but received an error %v", da, got)
	}
}

func TestDateAttributeValidateInvalid(t *testing.T) {
	da := DateAttribute("another_date")
	got := da.Validate()
	if got == nil {
		t.Errorf("expected an error, but received an error %v", got)
	}
}
