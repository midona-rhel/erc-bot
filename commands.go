package main

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

func commandHandler(session disgord.Session, evt *disgord.MessageCreate) error {
	if strings.HasPrefix(evt.Message.Content, config.Prefix) && !evt.Message.Author.Bot {

	}
	return nil
}

var commands = []command{
	{
		identifier: "help",
		name:       "help",
		f:          helpCommand,
	},
}

func helpCommand(session disgord.Session, evt *disgord.MessageCreate) error {
	m, err := evt.Message.Reply(evt.Ctx, session, config.Strings.HelpMessage)
	if err != nil {
		return fmt.Errorf("Failed to execute help command %w", err)
	} else if !evt.Message.IsDirectMessage() {
		cleanup(session, config.MessageShortDuration, m, evt.Message)
	}
	return nil
}

func cleanup(sessiom disgord.Session, duration int, ms ...*disgord.Message) error {

}

type command struct {
	identifier string
	name       string
	f          func(session disgord.Session, evt *disgord.MessageCreate) error
}
