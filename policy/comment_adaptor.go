package policy

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"regexp"
)

// CommitCommentEventAdaptor wraps the gitlab.CommitCommentEvent
type CommitCommentEventAdaptor struct {
	gitlab.CommitCommentEvent
}

// State provides access to the event's state if it has one
func (c CommitCommentEventAdaptor) State() (*string, error) {
	return nil, fmt.Errorf("commit comment events have no state")
}

// ResourceType returns the webhook's X-GitLab header value as an EventType
func (c CommitCommentEventAdaptor) ResourceType() gitlab.EventType {
	return gitlab.EventType(c.ObjectKind)
}

// Milestone returns the webhook's ID
func (c CommitCommentEventAdaptor) Milestone() (*int, error) {
	return nil, fmt.Errorf("commit comment events have no milestones")
}

// Labels returns the labels for a MergeEvent
func (c CommitCommentEventAdaptor) Labels() ([]string, error) {
	return nil, fmt.Errorf("commit comment events have no labels")
}

func (c CommitCommentEventAdaptor) Note() (*string, error) {
	return &c.ObjectAttributes.Note, nil
}

func (c CommitCommentEventAdaptor) Mentions() []string {
	re := regexp.MustCompile(`\@[a-zA-Z\.\-0-9]+`)
	mentions := re.FindAllStringSubmatch(c.ObjectAttributes.Note, -1)
	var people []string
	for _, p := range mentions {
		people = append(people, p...)
	}
	return people
}

func (c CommitCommentEventAdaptor) NoteType() (*string, error) {
	return &c.ObjectAttributes.NoteableType, nil
}
