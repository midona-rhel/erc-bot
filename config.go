package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// readconfig returns the config for the bot and panics if the read fails.
func loadConfig() Config {
	f, err := os.Open("config.json")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	c := Config{}
	if json.Unmarshal(b, &c) != nil {
		panic(err)
	}
	return c
}

// Config represents the runtime configuration saved as a json.
type Config struct {
	Token  string
	Prefix string

	MessageShortDuration int
	MessageLongDuration  int

	Strings struct {
		HelpMessage string
	}
}
