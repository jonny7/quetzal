package policy

import "fmt"

// DateAttribute is the updated or created property
type DateAttribute string

// DateCondition is the greater than or less than [date] filter
type DateCondition string

// DateIntervalType is the type of available interval
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
