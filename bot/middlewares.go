package bot

import (
	"github.com/go-chi/render"
	"net/http"
)

// webhookSecret confirms that a webhook secret matched the preconfigured one
// on the bot. If no secret exists in GitLab it's initialized with its zero value
func (b *Bot) webhookSecret(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Gitlab-Token") == "" {
			next.ServeHTTP(w, r)
		}
		if r.Header.Get("X-Gitlab-Token") != b.Config.Secret {
			w.WriteHeader(400)
			render.Respond(w, r, "webhook secret mismatch")
			return
		}
		next.ServeHTTP(w, r)
	}
}
