package policy

import (
	"github.com/xanzy/go-gitlab"
)

// Webhook is a wrapper around the incoming webhook
type Webhook struct {
	// EventType is the type of webhook, taken from the X-GitLab-Header
	EventType gitlab.EventType
	// Event is the payload of the webhook
	Event interface{}
}

type WebhookResult struct {
	Policy Policy
	Empty  bool
}

// FilterEvent processing policies against the incoming hook and only returns policies
// that are valid for this webhook Event.
func (w *Webhook) FilterEvent(in <-chan Policy) <-chan WebhookResult {
	validPolicies := make(chan WebhookResult)
	go func() {
		for pol := range in {
			switch ev := w.Event.(type) {
			case gitlab.CommitCommentEvent:
				cce := CommitCommentEventAdaptor{ev}
				validPolicies <- <-pol.ConditionsMet(cce)
			case gitlab.MergeEvent:
				me := MergeEventAdaptor{ev}
				validPolicies <- <-pol.ConditionsMet(me)
			// @todo these fail to be decoded when using the payload from GitLab docs
			case gitlab.MergeCommentEvent:
			case gitlab.SnippetCommentEvent:
			case gitlab.WikiPageEvent:
				we := WikiEventAdaptor{ev}
				validPolicies <- <-pol.ConditionsMet(we)
			}
		}
		close(validPolicies)
	}()
	return validPolicies
}
