package bot

import (
	"encoding/json"
	"io"
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

// decodeWebhook decodes webhook from Gitlab
func decodeWebhook(body io.Reader) (*Webhook, error) {
	var webhook Webhook
	err := json.NewDecoder(body).Decode(&webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (w *Webhook) handleEvent(options Config) (interface{}, error) {
	return nil, nil
}
