package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

type command struct {
	identifier string
	name       string
	perms      disgord.PermissionBits
	f          func(session disgord.Session, evt *disgord.MessageCreate, args string) (string, error)
}

var commands = []command{
	{
		identifier: "help",
		name:       "help",
		perms:      disgord.PermissionSendMessages,
		f:          helpCommand,
	},
}

func helpCommand(session disgord.Session, evt *disgord.MessageCreate, args string) (string, error) {
	m, err := evt.Message.Reply(evt.Ctx, session, config.CommandService.Strings.HelpMessage)
	if err != nil {
		return "", fmt.Errorf("Failed to send help message; %w", err)
	} else if !evt.Message.IsDirectMessage() {
		go cleanup(session, config.MessageShortDuration, m, evt.Message)
	}
	return "Help message sent", nil
}
