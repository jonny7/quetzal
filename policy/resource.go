package policy

import (
	"fmt"
	"strings"
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

// Validate ensures that an invalid EventType can't be used by the user provided policy
func (e EventType) Validate() error {
	switch e {
	case Push, Tag, Issue, Comment, MergeRequest, Wiki, Pipeline, Job, Deployment, Default:
		return nil
	}
	return fmt.Errorf("the inputted EventType is invalid: %s", e)
}

func (e EventType) ToLower() string {
	return strings.ToLower(string(e))
}
