package policy

import (
	"github.com/xanzy/go-gitlab"
)

type Webhook struct {
	EventType gitlab.EventType
	Event     interface{}
}

// FilterEvent processing policies against the incoming hook and only returns policies
// that are valid for this webhook Event.
func (w *Webhook) FilterEvent(in <-chan Policy) <-chan Policy {
	validPolicies := make(chan Policy)
	go func() {
		for pol := range in {
			switch ev := w.Event.(type) {
			case gitlab.MergeEvent:
				me := MergeEventAdaptor{ev}
				validPolicies <- <-pol.ConditionsMet(me)
			case gitlab.CommitCommentEvent:
				//if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				//	validPolicies <- <-pol.ConditionsMet()
				//}
			// @todo these fail to be decoded when using the payload from GitLab docs
			// case gitlab.MergeCommentEvent:
			// case gitlab.SnippetCommentEvent:
			case gitlab.IssueCommentEvent:
				//if pol.Conditions.Note.Type == nil || ev.ObjectAttributes.NoteableType == pol.Conditions.Note.Type.ToString() {
				//	validPolicies <- <-pol.ConditionsMet()
				//}
				//default:
				//	validPolicies <- <-pol.ConditionsMet()
			}
		}
		close(validPolicies)
	}()
	return validPolicies
}
