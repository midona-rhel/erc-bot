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
	config    *Config
	lastposts *UserLastPost
}

func main() {
	config := readconfig()
	d, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		panic(err)
	}
	b := new(Bot)
	b.lastposts = newUserLastPost()

	d.AddHandler(b.handleRoles)
	d.AddHandler(b.handleThrottle)

	if err = d.Open(); err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("ERC-BOT is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
