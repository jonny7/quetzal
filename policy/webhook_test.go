package policy

import (
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func readerToPolicyChan(p string) <-chan Policy {
	body, err := ioutil.ReadAll(io.NopCloser(strings.NewReader(p)))
	if err != nil {
		log.Fatalf("failed to load policies: %v", err)
	}

	var policies Policies
	err = yaml.Unmarshal(body, &policies)
	if err != nil {
		log.Fatalf("failed to unmarshal policies: %v", err)
	}

	out := make(chan Policy)
	go func() {
		defer close(out)
		for _, ruleSet := range policies.Policies {
			out <- ruleSet
		}
	}()
	return out
}

func TestWebhook(t *testing.T) {
	hook := gitlab.MergeEvent{}
	hook.ObjectAttributes.State = string(mergeRequestStateOpen)
	hook.Labels = []*gitlab.Label{{Name: "api"}, {Name: "critical"}}
	webhook := Webhook{EventType: gitlab.EventTypeMergeRequest, Event: hook}

	p := `policies:
  - name: Assign Critical Merges to Snr Staff
    resource: Merge Request Hook
    conditions:
      labels:
        - critical`

	policyChan := readerToPolicyChan(p)

	var results []WebhookResult
	got := webhook.FilterEvent(policyChan)

	for r := range got {
		results = append(results, r)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 webhook results, but got: %d", len(results))
	}
	if results[0].Empty != false {
		t.Errorf("expected Webhook Result to be Empty as hook must contain all labels listed on policy")
	}
}
