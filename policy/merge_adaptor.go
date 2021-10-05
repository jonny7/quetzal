package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

// MergeEventAdaptor wraps the gitlab.MergeEvent
type MergeEventAdaptor struct {
	gitlab.MergeEvent
}

// State provides access to the event's state if it has one
func (m MergeEventAdaptor) State() (*string, error) {
	return &m.ObjectAttributes.State, nil
}

// ResourceType returns the webhook's X-GitLab header value as an EventType
func (m MergeEventAdaptor) ResourceType() gitlab.EventType {
	return gitlab.EventType(m.ObjectKind)
}

// Milestone returns the webhook's ID
func (m MergeEventAdaptor) Milestone() (*int, error) {
	return &m.ObjectAttributes.MilestoneID, nil
}

// Labels returns the labels for a MergeEvent
func (m MergeEventAdaptor) Labels() ([]string, error) {
	var labels []string
	for _, label := range m.MergeEvent.Labels {
		labels = append(labels, label.Name)
	}
	return labels, nil
}

func (m MergeEventAdaptor) Note() (*string, error) {
	return nil, fmt.Errorf("MergeEvent hooks don't have a note field")
}

func (m MergeEventAdaptor) Mentions() []string {
	return nil
}

func (m MergeEventAdaptor) NoteType() (*string, error) {
	return nil, fmt.Errorf("MergeEvent hooks don't have a note type")
}
