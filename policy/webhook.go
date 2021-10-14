package policy

import (
	"github.com/xanzy/go-gitlab"
	"sort"
)

// Webhook is a wrapper around the incoming webhook
type Webhook struct {
	// EventType is the type of webhook, taken from the X-GitLab-Header
	EventType gitlab.EventType
	// Event is the payload of the webhook
	Event interface{}
}

// WebhookResult returns the Policy that was matched against the Webhook
// along with an array of actions that the bot tried to take.
type WebhookResult struct {
	Policy  Policy               `json:"policy"`
	Actions []GitLabUpdateResult `json:"actions"`
}

// FilterEvent takes a channel of Policy to check against the incoming Webhook
// If a Policy has conditions that are met by the Webhook, the Policy Action(s)
// are triggered, which normally results in the gitlab.Client making updates to GitLab
func (w *Webhook) FilterEvent(in <-chan Policy, client *gitlab.Client) <-chan WebhookResult {
	processed := make(chan WebhookResult)
	go func() {
		for pol := range in {
			result := WebhookResult{Policy: pol}
			switch ev := w.Event.(type) {
			case *gitlab.MergeEvent:
				me := MergeEventAdaptor{*ev}
				if matcher(pol, me, w.EventType) {
					result.Actions = me.execute(pol.Actions, client)
				}
				processed <- result
			}
		}
		close(processed)
	}()
	return processed
}

// matcher compares the Policy with the Webhook Adaptor and ultimately decides
// whether to process the Action
func matcher(policy Matcher, adaptor GitLabAdaptor, event gitlab.EventType) bool {
	if policy.resource() != event {
		return false
	}
	if policy.state() != nil {
		if !sliceContains(policy.state(), adaptor.state()[0]) {
			return false
		}
	}
	if policy.labels() != nil {
		sort.Slice(policy.labels(), func(i, j int) bool {
			return policy.labels()[i] < policy.labels()[j]
		})
		adaptorLabels := adaptor.labels()
		sort.Slice(adaptorLabels, func(i, j int) bool {
			return adaptorLabels[i] < adaptorLabels[j]
		})
		if !slicesMatch(policy.labels(), adaptorLabels) {
			return false
		}
	}
	return true
}

// checks if two slices match exactly. Expects the slices to have been sorted
func slicesMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sliceContains(a []string, b string) bool {
	for _, elem := range a {
		if b == elem {
			return true
		}
	}
	return false
}
