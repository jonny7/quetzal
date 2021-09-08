package bot

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/jonny7/quetzal/policy"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Config is the user declared details provided from the yaml file
// it contains general info for the bot along with a `Policies` property
type Config struct {
	User       string          `yaml:"user"`
	Token      string          `yaml:"token"`
	RepoHost   string          `yaml:"repoHost"`
	BotServer  string          `yaml:"botServer"`
	Endpoint   string          `yaml:"endpoint"`
	Secret     string          `yaml:"secret"`
	Port       string          `yaml:"port"`
	policyPath string          `yaml:"policyPath"`
	Policies   []policy.Policy `yaml:"policies"`
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
	r.Get("/policies", b.policies())
	r.Post("/reload", b.reload())
}

// processWebhook is the main endpoint for this bot
// and bridges user specified plugins and the gitlab webhook
func (b *Bot) processWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webhook, err := decodeWebhook(r.Body)
		if err != nil {
			render.Respond(w, r, Message{Msg: fmt.Sprintf("Could not decode webhook: %v", err)})
			return
		}
		// @todo complete this handler, omitted for small unit tests
		_, err = webhook.handleEvent(*b.Config)
		if err != nil {
			render.Respond(w, r, Message{Msg: fmt.Sprintf("Some error occurred: %v", err)})
			return
		}
		render.Respond(w, r, nil)
	}
}

// ping just provides a simple bot health endpoint
func (b *Bot) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, Message{Msg: "pong"})
	}
}

// policies endpoint returns all the loaded policies for the bot
func (b *Bot) policies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, b.Config.Policies)
	}
}

// reload will attempt to reload the bot's policies
func (b *Bot) reload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reader, err := createReader(b.Config.policyPath)
		if err != nil {
			w.WriteHeader(500)
			render.Respond(w, r, Message{Msg: fmt.Sprintf("could not create reader for policy file: %v", err)})
			return
		}
		err = b.loadPolicies(reader)
		if err != nil {
			w.WriteHeader(500)
			render.Respond(w, r, Message{Msg: fmt.Sprintf("policies couldn't be reloaded: %v", err)})
			return
		}
		render.Respond(w, r, Message{Msg: "policies reloaded"})
	}
}

// createReader takes a file location and returns an io.ReaderCloser
func createReader(file string) (io.ReadCloser, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return f, err
}

// loadConfig takes an io.ReaderCloser
// and attempts to load it into the Config struct
func loadConfig(reader io.ReadCloser) (*Config, error) {
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Error().Msg(fmt.Sprintf("the config file failed to close: %v", err))
		}
	}(reader)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var config Config
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("config file could not be read, error:%v", err))
		return nil, err
	}

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("config file could not be unmarshalled, error:%v", err))
		return nil, err
	}
	return &config, nil
}

// loadPolicies loads the specified policies.yml file
func (b *Bot) loadPolicies(reader io.ReadCloser) error {
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			b.Logger.Error().Msg(fmt.Sprintf("the config file failed to close: %v", err))
		}
	}(reader)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("policy file could not be read, error:%v", err))
		return err
	}

	var policies policy.Policies
	err = yaml.Unmarshal(body, &policies)
	if err != nil {
		b.Logger.Info().Msg(fmt.Sprintf("policy file could not be unmarshalled, error:%v", err))
		return err
	}
	b.Config.Policies = policies.Policies
	return nil
}

// New creates a new bot taking the config filename and path from `main`'s arguments
func New(config, policies string) (*Bot, error) {
	logger := zerolog.New(os.Stdout)
	reader, err := createReader(config)
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("an error occured creating a reader for the config file: %v", err))
	}
	c, err := loadConfig(reader)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		Router: chi.NewRouter(),
		Logger: logger,
		Config: c,
	}

	p, err := createReader(policies)
	if err != nil {
		b.Logger.Error().Msg(fmt.Sprintf("an error occured creating a reader for the policy file: %v", err))
	}
	err = b.loadPolicies(p)
	if err != nil {
		b.Logger.Error().Msg(fmt.Sprintf("policies couldn't be loaded: %v", err))
	}

	b.Router.Use(render.SetContentType(render.ContentTypeJSON))
	b.routes(b.Router)
	return b, nil
}
