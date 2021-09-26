package policy

import (
	"github.com/xanzy/go-gitlab"
)

// Stater returns the state of an event if possible
type Stater interface {
	State() (string, error)
}

// Resourcer returns the type of resource
type Resourcer interface {
	ResourceType() gitlab.EventType
}

// GitLabAdaptor wraps all the events
type GitLabAdaptor interface {
	Stater
	Resourcer
}
