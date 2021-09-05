package bot

import (
	"fmt"
	"os"
)

// EventType represents a GitLab Webhook type as string
type EventType string

// Available EventTypes
const (
	Push         EventType = "push"
	Tag          EventType = "tag_push"
	Issue        EventType = "issue"
	Comment      EventType = "note"
	MergeRequest EventType = "merge_request"
	Wiki         EventType = "wiki_page"
	Pipeline     EventType = "pipeline"
	Job          EventType = "build"
	Deployment   EventType = "deployment"
	Default      EventType = "default"
)

// Webhook is a minimal representation of GitLab's webhook
// it simply decodes in a type-safe way the event type
type Webhook struct {
	ObjectKind EventType `json:"object_kind"`
}

func (w *Webhook) handleEvent(options Config) (interface{}, error) {
	return nil, nil
}

func (w *Webhook) source(name string) ([]byte, error) {
	f, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %v", err)
	}
	return f, nil
}
