package policy

import "fmt"

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

const (
	createdAt DateAttribute = "created_at"
	updatedAt DateAttribute = "updated_at"
)

// validate confirms only the values created_at and updated_at were input
func (d *DateAttribute) validate() error {
	switch *d {
	case createdAt, updatedAt:
		return nil
	}
	return fmt.Errorf("`date:attribute` allowed options are: `%s`, `%s`, But received: %v", createdAt, updatedAt, d)
}

// DateCondition is the greater than or less than [date] filter
type DateCondition string

const (
	olderThan DateCondition = "older_than"
	newerThan DateCondition = "newer_than"
)

// validate confirms that only older_than and newer_than are passed into the config
func (d *DateCondition) validate() error {
	switch *d {
	case olderThan, newerThan:
		return nil
	}
	return fmt.Errorf("`date:condition` allowed options are: `%s`, `%s`. But received: %v", olderThan, newerThan, d)
}

// DateIntervalType is the type of available interval
type DateIntervalType string

const (
	days   DateIntervalType = "days"
	weeks  DateIntervalType = "weeks"
	months DateIntervalType = "months"
	years  DateIntervalType = "years"
)

// validate confirms that only days, weeks, months, years
// are passed in
func (d *DateIntervalType) validate() error {
	switch *d {
	case days, weeks, months, years:
		return nil
	}
	return fmt.Errorf("`date:intervalType` allowed options are: `%s`, `%s`, `%s`, `%s`. But received: %v", days, weeks, months, years, d)
}
