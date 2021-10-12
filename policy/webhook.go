package policy

import "github.com/xanzy/go-gitlab"

// Webhook is a wrapper around the incoming webhook
type Webhook struct {
	// EventType is the type of webhook, taken from the X-GitLab-Header
	EventType gitlab.EventType
	// Event is the payload of the webhook
	Event interface{}
}

type WebhookResult struct {
	Policy  Policy               `json:"policy"`
	Actions []GitLabUpdateResult `json:"actions"`
}

func (w *Webhook) FilterEvent(in <-chan Policy, client *gitlab.Client) <-chan WebhookResult {
	processed := make(chan WebhookResult)
	go func() {
		for pol := range in {
			result := WebhookResult{Policy: pol}
			switch ev := w.Event.(type) {
			case *gitlab.MergeEvent:
				me := MergeEventAdaptor{*ev}
				if pol.matcher(w.EventType, me) {
					result.Actions = me.execute(pol.Actions, client)
				}
				processed <- result
			}
		}
		close(processed)
	}()
	return processed
}
