package main

import (
	"fmt"
	"gitlab.com/jonny7/quetzal/bot"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	user := getEnvStr("user", "")
	token := getEnvStr("token", "")
	botServer := getEnvStr("bot", "")
	endpoint := getEnvStr("webhook", "/webhook-endpoint")
	secret := getEnvStr("secret", "")
	host := getEnvStr("host", "https://gitlab.com")
	port := getEnvStr("port", "7838")
	dry, err := getEnvBool("dry", false)
	if err != nil {
		log.Fatalf("dry var was unprocessible: %v", err)
	}
	policies := getEnvStr("policies", "./examples/.policies.yaml")
	version, err := getEnvBool("version", false)
	if err != nil {
		log.Fatalf("version var was unprocessible: %v", err)
	}

	if version {
		fmt.Println("Quetzal version ", getVersion())
		os.Exit(0)
	}

	config := bot.Config{
		User:       user,
		Token:      token,
		BotServer:  botServer,
		Endpoint:   endpoint,
		Secret:     secret,
		Host:       host,
		Port:       fmt.Sprintf(":%s", port),
		PolicyPath: policies,
		DryRun:     dry,
	}

	errorCh := make(chan error)
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		errorCh <- fmt.Errorf("%s", <-sigs)
	}()

	go func() {
		errorCh <- run(config, policies)
	}()

	log.Printf("exiting %v", <-errorCh)
}

func getEnvBool(key string, defaultVal bool) (bool, error) {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultVal, err
		}
		return b, nil
	}
	return defaultVal, nil
}

func getEnvStr(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}

func run(config bot.Config, policies string) error {
	b, err := bot.New(config, policies)
	if err != nil {
		return err
	}
	if httpErr := http.ListenAndServe(b.Config.Port, b.Router); httpErr != nil {
		return httpErr
	}
	fmt.Printf("Quetzal running with the following configuration: %v", b.Config)
	return nil
}
