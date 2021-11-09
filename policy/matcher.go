package policy

import "github.com/xanzy/go-gitlab"

// Matcher ensures that a Policy provides the required functionality
// to enable comparing a Webhook to a Policy
type Matcher interface {
	Stater
	Resourcer
	Labeler
	Milestoner
	ForbiddenLabeler
}

// Stater provides a method to get the state from an object
type Stater interface {
	state() []string
}

// Resourcer provides the resource type from a Policy
type Resourcer interface {
	resource() gitlab.EventType
}

// Labeler returns the labels for a Policy
type Labeler interface {
	labels() []string
}

// Milestoner provides the milestone ID
type Milestoner interface {
	milestone() int
}

type ForbiddenLabeler interface {
	forbiddenLabels() []string
}
