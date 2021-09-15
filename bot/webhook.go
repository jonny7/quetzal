package bot

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/jonny7/quetzal/policy"
)

type Webhook struct {
	eventType gitlab.EventType
	event     interface{}
}

// filterAdditionalEventType will further filter down the policies if there are additional sub-types
// the main example would be note events, which have 4 sub-types
func (w *Webhook) filterAdditionalEventType(policySubset []policy.Policy) ([]policy.Policy, []error) {
	var policies []policy.Policy
	var errors []error
	for _, p := range policySubset {
		switch w.event.(type) {
		case gitlab.CommitCommentEvent:
			cce, ok := w.event.(gitlab.CommitCommentEvent)
			if !ok {
				errors = append(errors, fmt.Errorf("type assertion for event %v of type %s failed", w.event, w.eventType))
				break
			}
			if p.Conditions.Note.Type == nil || p.Conditions.Note.Type.ToString() == cce.ObjectAttributes.NoteableType {
				policies = append(policies, p)
			}
		case gitlab.IssueCommentEvent:
			ice, ok := w.event.(gitlab.IssueCommentEvent)
			if !ok {
				errors = append(errors, fmt.Errorf("type assertion for event %v of type %s failed", w.event, w.eventType))
				break
			}

			if p.Conditions.Note.Type == nil || p.Conditions.Note.Type.ToString() == ice.ObjectAttributes.NoteableType {
				policies = append(policies, p)
			}
		default:
			policies = append(policies, p)
		}
	}

	return policies, errors
}
