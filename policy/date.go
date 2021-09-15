package policy

import "fmt"

// Validator provides a method to validate a struct
// created from a yaml config file, most types
// are backed by strings, which means that value can
// be passed in the .yml file. Validate confirms these are
// permitted values
type Validator interface {
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

const (
	createdAt DateAttribute = "created_at"
	updatedAt DateAttribute = "updated_at"
)

// Validate confirms only the values created_at and updated_at were input
func (d DateAttribute) Validate() error {
	switch d {
	case createdAt, updatedAt:
		return nil
	}
	return fmt.Errorf("`date:attribute` expected value of either `%s`, `%s`, But received: %v", d, createdAt, updatedAt)
}

// DateCondition is the greater than or less than [date] filter
type DateCondition string

const (
	olderThan DateCondition = "older_than"
	newerThan DateCondition = "newer_than"
)

// Validate confirms that only older_than and newer_than are passed into the config
func (d DateCondition) Validate() error {
	switch d {
	case olderThan, newerThan:
		return nil
	}
	return fmt.Errorf("`date:condition` expected values `%s`, `%s`. But received: %v", d, olderThan, newerThan)
}

// DateIntervalType is the type of available interval
type DateIntervalType string

const (
	days   DateIntervalType = "days"
	weeks  DateIntervalType = "weeks"
	months DateIntervalType = "months"
	years  DateIntervalType = "years"
)

// Validate confirms that only days, weeks, months, years
// are passed in
func (d DateIntervalType) Validate() error {
	switch d {
	case days, weeks, months, years:
		return nil
	}
	return fmt.Errorf("`date:intervalType` expected values `%s`, `%s`, `%s`, `%s`. But received: %v", d, days, weeks, months, years)
}
