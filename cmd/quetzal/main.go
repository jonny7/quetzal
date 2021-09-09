package main

import (
	"flag"
	"fmt"
	"gitlab.com/jonny7/quetzal/bot"
	"log"
	"net/http"
)

const v = "0.0.0"

func main() {
	var config, policies string

	flag.StringVar(&config, "config", "./config.yaml", "The relative path to the config file name and extension")
	flag.StringVar(&policies, "policies", "./.policies.yaml", "The relative path to the policies file")
	version := flag.Bool("version", false, "display version of quetzal")
	flag.Parse()

	if *version {
		fmt.Printf("Quetzal version %v", v)
	}

	if err := run(config, policies); err != nil {
		log.Fatalf("error launching bot %v", err)
	}
}

func run(config, policies string) error {
	b, err := bot.New(config, policies)
	if err != nil {
		return err
	}
	if err := http.ListenAndServe(b.Config.Port, b.Router); err != nil {
		return err
	}
	return nil
}
