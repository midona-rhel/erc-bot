package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// readconfig returns the config for the bot and panics if the read fails.
func readconfig() *Config {
	f, err := os.Open("config.json")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	c := new(Config)
	if json.Unmarshal(b, c) != nil {
		panic(err)
	}
	return c
}

// Config represents the runtime configuration saved as a json.
type Config struct {
	Discord struct {
		Token        string `json:"token"`
		AdminRole    string `json:"adminRole"`
		DefaultGuild string `json:"defaultGuild"`
		Playing      string `json:"playing"`
	} `json:"discord"`
	Role []struct {
		RoleID string   `json:"roleID"`
		Alias  []string `json:"alias"`
	} `json:"role"`
	Throttle []struct {
		ChannelID     string `json:"channelID"`
		MaxTokens     int    `json:"maxTokens"`
		TokenInterval int    `json:"tokenInterval"`
		CharLimit     int    `json:"charLimit"`
		NewlineLimit  int    `json:"newlineLimit"`
	} `json:"throttle"`
	Monitor struct {
		Output string `json:"output"`
	} `json:"monitor"`
	Purge []struct {
		ChannelID      string `json:"channelID"`
		CronExpression string `json:"cronExpression"`
	} `json:"purge"`
	Help           string `json:"help"`
	CommandPrefix  string `json:"commandPrefix"`
	WelcomeMessage string `json:"welcomeMessage"`
}
