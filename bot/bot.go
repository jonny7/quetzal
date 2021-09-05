package bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jonny7/gitlab-bot/policy"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
)

// Config is the user declared details provided from the yaml file
// it contains general info for the bot along with a `Policies` property
type Config struct {
	User      string          `yaml:"user"`
	Token     string          `yaml:"token"`
	RepoHost  string          `yaml:"repoHost"`
	BotServer string          `yaml:"botServer"`
	Endpoint  string          `yaml:"endpoint"`
	Secret    string          `yaml:"secret"`
	Port      string          `yaml:"port"`
	Policies  []policy.Policy `yaml:"policies"`
}

// Bot struct encapsulates all behaviour of the bot
type Bot struct {
	Router *chi.Mux
	Logger zerolog.Logger
	Config *Config
}

// Message provides a simple message struct for times you need some
// kind of json response, like checking the health of the bot
type Message struct {
	Msg string `json:"msg"`
}

// ServeHTTP implements the handler interface
func (b *Bot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Router.ServeHTTP(w, r)
}

// routes declares the routes that the bot responds to
func (b *Bot) routes(r *chi.Mux) {
	r.Post(b.Config.Endpoint, b.webhookSecret(b.processWebhook()))
	r.Get("/ping", b.ping())
}

// processWebhook is the main endpoint for this bot
// and bridges user specified plugins and the gitlab webhook
func (b *Bot) processWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webhook, err := decodeWebhook(r.Body)
		if err != nil {
			render.Respond(w, r, Message{Msg: fmt.Sprintf("Could not decode webhook: %v", err)})
		}
		_, err = webhook.handleEvent(*b.Config)
		if err != nil {
			render.Respond(w, r, Message{Msg: fmt.Sprintf("Some error occurred: %v", err)})
		}
		render.Respond(w, r, nil)
	}
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

// ping just provides a simple bot health endpoint
func (b *Bot) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, Message{Msg: "pong"})
	}
}

// loadConfig takes the user specified yaml file
// and loads it into the Config struct
func loadConfig(name, path string) (*Config, error) {
	f, err := os.ReadFile(path + name)
	if err != nil {
		return nil, fmt.Errorf("config file could not be loaded at location: %s%s", path, name)
	}
	var config Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return nil, fmt.Errorf("config file could not be unmarshalled, error:%v", err)
	}
	return &config, nil
}

// New creates a new bot taking the config filename and path from `main`'s arguments
func New(name, path string) (*Bot, error) {
	logger := zerolog.New(os.Stdout)
	config, err := loadConfig(name, path)
	if err != nil {
		return nil, err
	}

	s := &Bot{
		Router: chi.NewRouter(),
		Logger: logger,
		Config: config,
	}
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
	s.routes(s.Router)
	return s, nil
}
