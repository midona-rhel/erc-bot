package main

import "os"
import "encoding/json"
import "io/ioutil"

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
	} `json:"throttle"`
	Monitor struct {
		Output string   `json:"output"`
		Events []string `json:"events"`
	} `json:"monitor"`
	Welcome struct {
		Message string `json:"message"`
	} `json:"welcome"`
	Purge []struct {
		ChannelID      string `json:"channelID"`
		CronExpression string `json:"cronExpression"`
	} `json:"purge"`
	Accuracy struct {
	} `json:"accuracy"`
	Help struct {
	} `json:"help"`
	Clear struct {
	} `json:"clear"`
	CommandPrefix string `json:"commandPrefix"`
}
