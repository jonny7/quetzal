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
	policy  Policy
	actions []GitLabUpdateResult
}

func (w *Webhook) FilterEvent(in <-chan Policy, client *gitlab.Client) <-chan WebhookResult {
	processed := make(chan WebhookResult)
	go func() {
		for pol := range in {
			result := WebhookResult{policy: pol}
			switch ev := w.Event.(type) {
			case gitlab.MergeEvent:
				me := MergeEventAdaptor{ev}
				if pol.matcher(*w) {
					result.actions = me.execute(pol.Actions, client)
				}
				processed <- result
			}
		}
		close(processed)
	}()
	return processed
}
