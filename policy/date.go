package policy

import "fmt"

// Validate provides a method to validate a struct
// created from a yaml config file, most types
// are backed by strings, which means that value can
// be passed in the .yml file. Validate confirms these are
// permitted values
type Validate interface {
	Validate() error
}

// Date is possible condition that can be used to allow or
// disallow the behaviour of the Bot see `config.yaml`
type Date struct {
	// Attribute can be `created_at` or `updated_at`
	Attribute DateAttribute `yaml:"attribute"`
	// Condition can be `older_than` or `newer_than`
	Condition DateCondition `yaml:"condition"`
	// IntervalType can be `days`, `weeks`, `months`, `years`
	IntervalType DateIntervalType `yaml:"intervalType"`
	// Interval is a numeric representation of the `IntervalType`
	Interval int `yaml:"interval"`
}

// DateAttribute is the updated or created property
type DateAttribute string

// Validate confirms only the values created_at and updated_at were input
func (d DateAttribute) Validate() error {
	switch d {
	case "created_at",
		"updated_at":
		return nil
	}
	return fmt.Errorf("`date:attribute` expected values created_at, updated_at, But received: %v", d)
}

// DateCondition is the greater than or less than [date] filter
type DateCondition string

// Validate confirms that only older_than and newer_than are passed into the config
func (d DateCondition) Validate() error {
	switch d {
	case "older_than",
		"newer_than":
		return nil
	}
	return fmt.Errorf("`date:condition` expected values older_than, newer_than. But received: %v", d)
}

// DateIntervalType is the type of available interval
type DateIntervalType string

// Validate confirms that only days, weeks, months, years
// are passed in
func (d DateIntervalType) Validate() error {
	switch d {
	case "days",
		"weeks",
		"months",
		"years":
		return nil
	}
	return fmt.Errorf("`date:intervalType` expected values days, weeks, months, years. But received: %v", d)
}
