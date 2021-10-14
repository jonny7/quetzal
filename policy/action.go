package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

// Action struct houses how an eligible webhook
// event should be responded to
type Action struct {
	// HTTP contains allows a user to call an API endpoint
	// and use that in subsequent actions
	//HTTP *HTTP `yaml:"http,omitempty"` @todo
	// Labels identifies which labels to add to an issue
	Labels []string `yaml:"labels,omitempty"`
	// RemoveLabels defines what labels to remove
	RemoveLabels []string `yaml:"removeLabels,omitempty"`
	// Status sets the status of the issue
	Status string `yaml:"status,omitempty"`
	// Mention allows the bot to mention certain users
	Mention []string `yaml:"mention,omitempty"`
	// Comment will leave a comment on said issue
	Comment string `yaml:"comment,omitempty"`
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
	if a.RemoveLabels != nil || a.Labels != nil {
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
	if a.Status != "" {
		return true
	}
	return false
}

// validate the actions are possible based on the webhook
func (a Action) validate(eventType gitlab.EventType) error {
	// validate status action against the type of webhook it wants to update
	if a.Status == "" {
		return nil
	}
	state := &State{State: a.Status}
	if err := state.validate(eventType); err != nil {
		return err
	}
	return nil
}
