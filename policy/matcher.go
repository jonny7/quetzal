package policy

// Matcher ensures that a Policy provides the required functionality
// to enable validation and matching
type Matcher interface {
	Stater
}
