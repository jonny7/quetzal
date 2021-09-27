package bot

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/jonny7/quetzal/policy"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

// Config is the user declared details provided from the yaml file
// it contains general info for the bot along with a `Policies` property
type Config struct {
	User       string
	Token      string
	BotServer  string
	Endpoint   string
	Secret     string
	Port       string
	PolicyPath string
	DryRun     bool
	Policies   []policy.Policy `yaml:"policies"`
}

// Bot struct encapsulates all behaviour of the bot
type Bot struct {
	Router *chi.Mux
	Logger *zerolog.Logger
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

// policies endpoint returns all the loaded policies for the bot
func (b *Bot) policies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, b.Config.Policies)
	}
}

// reload will attempt to reload the bot's policies
func (b *Bot) reload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reader, err := createReader(b.Config.PolicyPath)
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

// ping just provides a simple bot health endpoint
func (b *Bot) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, Message{Msg: "pong"})
	}
}

// processWebhook is the main endpoint for this bot
// and bridges user specified plugins and the gitlab webhook
func (b *Bot) processWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			render.Respond(w, r, nil)
			return
		}
		eventType := gitlab.HookEventType(r)
		event, err := gitlab.ParseWebhook(eventType, payload)
		if err != nil {
			render.Respond(w, r, Message{Msg: fmt.Sprintf("Could not decode webhook: %v", err)})
			return
		}
		webhook := policy.Webhook{
			EventType: eventType,
			Event:     event,
		}

		preparedPolicies := b.preparePolicies()
		workers := make([]<-chan policy.Policy, runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			workers[i] = webhook.FilterEvent(preparedPolicies)
		}

		validPolicies := merge(workers...)
		var p []policy.Policy
		for v := range validPolicies {
			p = append(p, v)
		}

		render.Respond(w, r, p)
	}
}

func merge(incoming ...<-chan policy.Policy) <-chan policy.Policy {
	var wg sync.WaitGroup

	wg.Add(len(incoming))
	outgoing := make(chan policy.Policy)
	multiplexer := func(policies <-chan policy.Policy) {
		defer wg.Done()
		for p := range policies {
			outgoing <- p
		}
	}

	for _, ch := range incoming {
		go multiplexer(ch)
	}

	go func() {
		wg.Wait()
		close(outgoing)
	}()

	return outgoing
}

// preparePolicies creates a channel of policies to process
// concurrently.
func (b *Bot) preparePolicies() <-chan policy.Policy {
	out := make(chan policy.Policy)
	go func() {
		defer close(out)
		for _, ruleSet := range b.Config.Policies {
			out <- ruleSet
		}
	}()
	return out
}

// createReader takes a file location and returns an io.ReaderCloser
func createReader(file string) (io.ReadCloser, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return f, err
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

// validatePolicies validates all the policies and fields where only certain values are allowed
// @todo finish validations
func (b *Bot) validatePolicies() error {
	policiesToValidate := b.preparePolicies()
	for p := range policiesToValidate {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("policy named: %s failed validation: %v", p.Name, err)
		}
	}
	return nil
}

// New creates a new bot taking the config filename and path from `main`'s arguments
func New(config Config, policies string) (*Bot, error) {
	logger := zerolog.New(os.Stdout)

	b := &Bot{
		Router: chi.NewRouter(),
		Logger: &logger,
		Config: &config,
	}

	p, err := createReader(policies)
	if err != nil {
		b.Logger.Error().Msg(fmt.Sprintf("an error occurred creating a reader for the policy file: %v", err))
	}
	if err = b.loadPolicies(p); err != nil {
		b.Logger.Error().Msg(fmt.Sprintf("policies couldn't be loaded: %v", err))
	}
	if err = b.validatePolicies(); err != nil {
		b.Logger.Error().Msg(fmt.Sprintf("invalid policy: %v", err))
	}

	b.Router.Use(render.SetContentType(render.ContentTypeJSON))
	b.Router.Use(middleware.Recoverer)
	b.routes(b.Router)
	return b, nil
}
