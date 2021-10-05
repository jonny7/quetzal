package policy

import (
	"github.com/xanzy/go-gitlab"
)

// Stater returns the state of an event if possible
type Stater interface {
	State() (*string, error)
}

// Resourcer returns the type of resource
type Resourcer interface {
	ResourceType() gitlab.EventType
}

// Milestoner returns the milestone if possible or error
type Milestoner interface {
	Milestone() (*int, error)
}

// Labeler returns the labels if available
type Labeler interface {
	Labels() ([]string, error)
}

// Noter returns information about the Note and NoteType
type Noter interface {
	Note() (*string, error)
	Mentions() []string
	NoteType() (*string, error)
}

// GitLabAdaptor wraps all the events
type GitLabAdaptor interface {
	Stater
	Resourcer
	Milestoner
	Labeler
	Noter
}
