package policy

import "fmt"

type DateAttribute string
type DateCondition string
type DateIntervalType string

// Date is possible condition that can be used to allow or
// disallow the behaviour of the Bot see `config.yaml`
type Date struct {
	// Attribute can be `created_at` or `updated_at`
	Attribute DateAttribute `yaml:"attribute"`
	// Condition can be `older_than` or `newer_than`
	Condition DateCondition `yaml:"condition"`
	// IntervalType can be `days`, `weeks`, `months`, `years`
	IntervalType DateIntervalType `yaml:"interval_type"`
	// Interval is a numeric representation of the `IntervalType`
	Interval int `yaml:"interval"`
}

func (dit DateIntervalType) validate() error {
	switch dit {
	case "days",
		"weeks",
		"months",
		"years":
		return nil
	}
	return fmt.Errorf("expected values days, weeks, months, years. But received: %v", dit)
}
