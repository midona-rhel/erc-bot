package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the runtime instance of the erc-bot.
type Bot struct {
	session           *discordgo.Session
	config            *Config
	throttledChannels *throttledChannelUserTokenMap
}

func main() {
	initLog()
	config := readconfig()
	initMonitor()
	session, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Panic(err)
	}
	bot := new(Bot)
	bot.config = config
	bot.throttledChannels = newThrottledChannelUserTokenMap()

	// Add handlers
	session.AddHandler(bot.handleCommands)
	session.AddHandler(bot.handleThrottle)
	// Add monitors
	session.AddHandler(bot.monitorGuildAdd)
	session.AddHandler(bot.monitorGuildRemove)
	session.AddHandler(bot.monitorMessageCreate)
	session.AddHandler(bot.monitorMessageDelete)
	session.AddHandler(bot.monitorMessageUpdate)

	bot.session = session

	if err = session.Open(); err != nil {
		log.Panic(err)
	}

	bot.purge(session)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("ERC-BOT is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
