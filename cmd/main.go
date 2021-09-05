package main

import (
	"flag"
	"gitlab.com/jonny7/quetzal/bot"
	"log"
	"net/http"
)

func main() {
	var name, path string

	flag.StringVar(&name, "config", "config.yaml", "The config file name (no extension)")
	flag.StringVar(&path, "path", "./", "The path to the config file, including directory separator")
	flag.Parse()

	if err := run(name, path); err != nil {
		log.Fatalf("error launching bot %v", err)
	}
}

func run(name, path string) error {
	b, err := bot.New(name, path)
	if err != nil {
		return err
	}
	if err := http.ListenAndServe(b.Config.Port, b.Router); err != nil {
		return err
	}
	return nil
}
