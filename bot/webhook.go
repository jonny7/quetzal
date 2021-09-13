package bot

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jonny7/quetzal/policy"
	"io"
)

// Webhook is a minimal representation of GitLab's webhook
// it simply decodes in a type-safe way the event type
type Webhook struct {
	ObjectKind policy.EventType `json:"object_kind"`
}

// decodeWebhook decodes webhook from Gitlab
func decodeWebhook(body io.Reader) (*Webhook, error) {
	var webhook Webhook
	err := json.NewDecoder(body).Decode(&webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (w *Webhook) handleEvent(bot *Bot) (interface{}, error) {
	matchedPolicies := bot.triggeredPolicies()
	if len(matchedPolicies) < 1 {
		bot.Logger.Info().Msg(fmt.Sprintf("no policies matched for this event: %v", w))
		return nil, nil
	}
	if bot.Config.DryRun {
		bot.Logger.Info().Msg(fmt.Sprintf("dry-run is true: so returning policies: %v", matchedPolicies))
		return nil, nil
	}
	return nil, nil
}
