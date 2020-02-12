package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

// The CommandService serves commands and logs them
type CommandService struct {
	logger Logger
}

func (c *CommandService) handler(session disgord.Session, evt *disgord.MessageCreate) error {
	if strings.HasPrefix(evt.Message.Content, config.CommandService.Prefix) && !evt.Message.Author.Bot {
		s := strings.Replace(evt.Message.Content, config.CommandService.Prefix, "", 1)
		s = strings.TrimSpace(s)
		for _, c := range commands {

			if !strings.HasPrefix(s, c.identifier) {
				continue
			}
			s = strings.Replace(s, c.identifier, "", 1)
			p, err := comparePermissions(session, evt, c.perms)

			if err != nil {
				logger.Error(fmt.Errorf("Failed to calculate command permissions; %w", err))
				break
			} else if p <= 0 {
				logger.Error(fmt.Errorf("Bot lacks permission to execute command: %s", c.name))
				break
			}

			s, err = c.f(session, evt, s)
			if err != nil {
				logger.Error(fmt.Errorf("Failed to execute command: %s; %w", c.name, err))
			} else {
				logger.Info(fmt.Sprintf("Executed command: %s; %s", c.name, s))
			}

			break
		}
	}
	return nil
}

func comparePermissions(session disgord.Session, evt *disgord.MessageCreate, perms disgord.PermissionBits) (uint64, error) {
	guildID := disgord.ParseSnowflakeString(config.GuildID)
	user, err := session.GetCurrentUser(context.Background())
	if err != nil {
		return 0, fmt.Errorf("Failed to get current user for permissions calculation; %w", err)
	}
	p, err := session.GetMemberPermissions(context.Background(), guildID, user.ID)
	if err != nil {
		return 0, fmt.Errorf("Failed to get permissions for current user; %w", err)
	}
	return p & perms, nil
}

var commands = []command{
	{
		identifier: "help",
		name:       "help",
		perms:      disgord.PermissionSendMessages,
		f:          helpCommand,
	},
}

func helpCommand(session disgord.Session, evt *disgord.MessageCreate, cleaned string) (string, error) {
	m, err := evt.Message.Reply(evt.Ctx, session, config.CommandService.Strings.HelpMessage)
	if err != nil {
		return "", fmt.Errorf("Failed to send help message; %w", err)
	} else if !evt.Message.IsDirectMessage() {
		cleanup(session, config.MessageShortDuration, m, evt.Message)
	}
	return "Help message sent", nil
}

func cleanup(sessiom disgord.Session, duration int, ms ...*disgord.Message) {

}

type command struct {
	identifier string
	name       string
	perms      disgord.PermissionBits
	f          func(session disgord.Session, evt *disgord.MessageCreate, cleaned string) (string, error)
}
