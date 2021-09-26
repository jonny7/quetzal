package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

// WikiEventAdaptor wraps the gitlab.WikiPageEvent
type WikiEventAdaptor struct {
	gitlab.WikiPageEvent
}

// State provides access to the event's state if it has one
func (w WikiEventAdaptor) State() (*string, error) {
	return nil, fmt.Errorf("WikiPageEvent has no `state`")
}

// ResourceType returns the webhook's X-GitLab header value as an EventType
func (w WikiEventAdaptor) ResourceType() gitlab.EventType {
	return gitlab.EventType(w.ObjectKind)
}

func (w WikiEventAdaptor) Milestone() (*int, error) {
	return nil, fmt.Errorf("WikiPageEvent has not Milestone")
}
