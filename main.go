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
	config            *Config
	throttledChannels *throttledChannelUserTokenMap
}

func main() {
	initLog()
	config := readconfig()
	initMonitor(config.Monitor.Output)
	session, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Panic(err)
	}
	bot := new(Bot)
	bot.config = config
	bot.throttledChannels = newThrottledChannelUserTokenMap()

	session.AddHandler(bot.handleCommands)
	session.AddHandler(bot.handleThrottle)

	session.AddHandler(monitorGuildAdd)
	session.AddHandler(monitorGuildRemove)
	session.AddHandler(monitorMessageCreate)
	session.AddHandler(monitorMessageDelete)
	session.AddHandler(monitorMessageUpdate)

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
