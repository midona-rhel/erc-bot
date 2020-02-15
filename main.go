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
	messages          *messageMap
}

func main() {
	initLog()
	log.Info("Starting v0.7.91")
	config := readconfig()
	initMonitor()
	session, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Panic(err)
	}
	session.StateEnabled = true
	session.State = discordgo.NewState()
	bot := new(Bot)
	bot.config = config
	bot.throttledChannels = newThrottledChannelUserTokenMap()
	bot.messages = &messageMap{
		messages: map[string]discordgo.MessageCreate{},
	}

	// Add handlers
	session.AddHandler(bot.handleCommands)
	session.AddHandler(bot.handleThrottle)
	session.AddHandler(bot.handleWelcomeMessage)
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
	bot.session.UpdateStatus(0, bot.config.Discord.Playing)

	bot.purge(session)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("ERC-BOT is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
