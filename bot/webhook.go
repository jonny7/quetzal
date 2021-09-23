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

// filterEvent processing policies against the incoming hook and only returns policies
// that are valid for this webhook event.
func (w *Webhook) filterEvent(in <-chan policy.Policy) <-chan policy.Policy {
	validPolicies := make(chan policy.Policy)
	go func() {
		for pol := range in {
			switch ev := w.event.(type) {
			case gitlab.CommitCommentEvent:
				//if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				fmt.Println(ev)
				validPolicies <- <-pol.ConditionsMet()
				//}
			// @todo these fail to be decoded when using the payload from GitLab docs
			// case gitlab.MergeCommentEvent:
			case gitlab.IssueCommentEvent:
				//if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				validPolicies <- <-pol.ConditionsMet()
				//}
			// @todo these fail to be decoded when using the payload from GitLab docs
			// case gitlab.SnippetCommentEvent:
			default:
				validPolicies <- <-pol.ConditionsMet()
			}
		}
		close(validPolicies)
	}()
	return validPolicies
}
