package policy

import "github.com/xanzy/go-gitlab"

// MergeEventAdaptor wraps the gitlab.MergeEvent
type MergeEventAdaptor struct {
	gitlab.MergeEvent
}

// State provides access to the event's state if it has one
func (m MergeEventAdaptor) State() (string, error) {
	return m.ObjectAttributes.State, nil
}

// ResourceType returns the webhook's X-GitLab header value as an EventType
func (m MergeEventAdaptor) ResourceType() gitlab.EventType {
	return gitlab.EventType(m.ObjectKind)
}
