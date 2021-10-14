package policy

import (
	"encoding/json"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"testing"
)

func TestWebhookFilter(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	oneUpdate := Policy{Name: "Test Webhook Filter", Resource: Resource{gitlab.EventTypeMergeRequest}, Actions: Action{Comment: "I am a bot"}}
	twoUpdates := Policy{Name: "2nd Test Webhook Filter", Resource: Resource{gitlab.EventTypeMergeRequest}, Actions: Action{Comment: "I am a bot", Labels: Label{[]string{"done"}}}}

	// mux response
	n := new(gitlab.Note)
	n.Body = oneUpdate.Actions.commentate()

	mergeEvent := gitlab.MergeEvent{}
	mergeEvent.Project.ID = 1
	mergeEvent.ObjectAttributes.IID = 234

	// mock response for add Note
	noteEndpoint := fmt.Sprintf("/api/v4/projects/%d/merge_requests/%d/notes", mergeEvent.Project.ID, mergeEvent.ObjectAttributes.IID)
	mux.HandleFunc(noteEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		noteErr := json.NewEncoder(w).Encode(&n)
		if noteErr != nil {
			t.Errorf("failed to encode response")
		}
		return
	})

	updateMergeRequestEndpoint := stubUpdatedMergeEventEndPoint(MergeEventAdaptor{mergeEvent})

	// response object for MergeRequest Updates, set to the action Labels
	m := new(gitlab.MergeRequest)
	m.Labels = twoUpdates.Actions.Labels.Labels

	// mock response for updateMerge Req
	mux.HandleFunc(updateMergeRequestEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&m)
		if err != nil {
			t.Errorf("failed to encode response")
		}
		return
	})

	data := []struct {
		name               string
		policy             Policy
		hook               Webhook
		expectedPolicyName string
		expectedActions    int
		expectedEndPoint   string
	}{
		{name: "Webhook that makes 1 update to GitLab", policy: oneUpdate, hook: Webhook{EventType: gitlab.EventTypeMergeRequest, Event: &mergeEvent}, expectedPolicyName: oneUpdate.Name, expectedActions: 1},
		{name: "Webhook that makes 2 update to GitLab", policy: twoUpdates, hook: Webhook{EventType: gitlab.EventTypeMergeRequest, Event: &mergeEvent}, expectedPolicyName: twoUpdates.Name, expectedActions: 2},
		// this uses a non-possible eventType, purely so the filter simulates a webhook and policy that match, but would fail
		// at some point in the conditions checks, as only one check is in place EventType, it manufactures this failure
		// to handle a nil result on the channel. This can be replaced when condition checks are built
		{name: "Webhook that is the correct type but policy, but doesn't match", policy: twoUpdates, hook: Webhook{EventType: gitlab.EventTypeIssue, Event: &mergeEvent}, expectedPolicyName: "", expectedActions: 0},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			ch := make(chan Policy)
			go func() {
				defer close(ch)
				ch <- d.policy
			}()
			got := <-d.hook.FilterEvent(ch, client)
			if got.Policy.Name != d.policy.Name {
				t.Errorf("expected policy name to be: %s", d.policy.Name)
			}

			if len(got.Actions) != d.expectedActions {
				t.Errorf("expected filter to return 1 updated function, but got %d", len(got.Actions))
			}
		})
	}
}
