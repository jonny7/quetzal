package main

import (
	"flag"
	"fmt"
	"gitlab.com/jonny7/quetzal/bot"
	"log"
	"net/http"
	"os"
)

func main() {
	var user, token, policies, botServer, endpoint, secret string
	var port int
	var dry, version bool

	flag.StringVar(&user, "user", "username@gitlab.com", "The Gitlab user this bot will act as")
	flag.StringVar(&token, "token", "notareatoken", "The personal access token for the stated user")
	flag.StringVar(&botServer, "bot-server", "https://bot-bot.com", "The base URL the bot lives on")
	flag.StringVar(&endpoint, "webhook-endpoint", "/webhook/path", "The webhook endpoint")
	flag.StringVar(&secret, "", "1234abcd", "The (optional) webhook secret ")
	flag.IntVar(&port, "port", 7838, "The port the bot listens on")
	flag.BoolVar(&dry, "dry-run", false, "don't perform any actions, just print out the actions that would be taken if live")
	flag.StringVar(&policies, "policies", "./.policies.yaml", "The relative path to the policies file")
	flag.BoolVar(&version, "version", false, "display version of quetzal")
	flag.Parse()

	if version {
		fmt.Println("Quetzal version ", current.toString())
		os.Exit(0)
	}

	config := bot.Config{
		User:       user,
		Token:      token,
		BotServer:  botServer,
		Endpoint:   endpoint,
		Secret:     secret,
		Port:       fmt.Sprintf(":%d", port),
		PolicyPath: policies,
		DryRun:     dry,
	}

	if err := run(config, policies); err != nil {
		log.Fatalf("error launching bot %v", err)
	}
}

func run(config bot.Config, policies string) error {
	b, err := bot.New(config, policies)
	if err != nil {
		return err
	}
	if httpErr := http.ListenAndServe(b.Config.Port, b.Router); httpErr != nil {
		return httpErr
	}
	return nil
}
