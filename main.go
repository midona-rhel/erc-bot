package main

import (
	"github.com/bwmarrin/discordgo"
)

// Bot represents the runtime instance of the erc-bot.
type Bot struct {
	config *Config
}

func main() {
	config := readconfig()
	discordgo.New("Bot" + config.Discord.Token)
}
