package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

// ActionStatus is a string backed status
type ActionStatus string

// Action struct houses how an eligible webhook
// event should be responded to
type Action struct {
	// Labels identifies which labels to add to an issue
	Labels Labels `yaml:",inline"`
	// RemoveLabels defines what labels to remove
	RemoveLabels []string `yaml:"removeLabels,omitempty"`
	// Status sets the status of the issue
	Status ActionStatus `yaml:"status"`
	// Mention allows the bot to mention certain users
	Mention []string `yaml:"mention,omitempty"`
	// Comment will leave a comment on said issue
	Comment string `yaml:"comment,omitempty"`
}

// Labels represents an array of labels
type Labels struct {
	Labels []string `yaml:"labels"`
}

// commentate builds a note body with the mentions and comment
// concatenated
func (a Action) commentate() string {
	var out string
	for _, m := range a.Mention {
		out += fmt.Sprintf("@%s ", m)
	}
	if len(a.Mention) == 0 {
		return a.Comment
	}
	return fmt.Sprintf("%s%s", out, a.Comment)
}

func (a Action) updateLabels() bool {
	if a.RemoveLabels != nil || a.Labels.Labels != nil {
		return true
	}
	return false
}

func (a Action) addNote() bool {
	if a.Mention != nil || a.Comment != "" {
		return true
	}
	return false
}

func (a Action) updateState() bool {
	return a.Status != ""
}

// validate the actions are possible based on the webhook
func (a Action) validate(eventType gitlab.EventType) error {
	// validate status action against the type of webhook it wants to update
	if a.Status == "" {
		return nil
	}
	if err := a.Status.fieldValidator(eventType); err != nil {
		return err
	}
	return nil
}

func (as ActionStatus) fieldValidator(eventType gitlab.EventType) error {
	if eventType == gitlab.EventTypeMergeRequest {
		switch strings.ToLower(string(as)) {
		case string(mergeRequestStateOpen), string(mergeRequestStateClose), string(mergeRequestStateApproved):
			return nil
		default:
			return fmt.Errorf("merge Request Events only allow for statuses of open, closed, approved")
		}
	}
	return fmt.Errorf("the status of %s for event type %s is invalid", as, eventType)
}
