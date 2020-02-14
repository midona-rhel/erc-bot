package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func commandHandler(session disgord.Session, evt *disgord.MessageCreate) {

	if strings.HasPrefix(evt.Message.Content, config.CommandService.Prefix) && !evt.Message.Author.Bot {

		s := strings.Replace(evt.Message.Content, config.CommandService.Prefix, "", 1)
		s = strings.TrimSpace(s)

		for _, c := range commands {

			if !strings.HasPrefix(s, c.identifier) {
				continue
			}
			s = strings.Replace(s, c.identifier, "", 1)

			p, err := evalPerms(session, c.perms)

			if err != nil {
				logger.Error(fmt.Errorf("Failed to execute command; %w", err))
				break
			} else if p == 0 {
				logger.Error(fmt.Errorf("Failed to execute command, bot lacks permissions: %s", c.name))
				break
			}

			s, err = c.f(session, evt, s)
			if err != nil {
				logger.Error(fmt.Errorf("Failed to execute command: %s; %w", c.name, err))
			} else {
				logger.Info(fmt.Sprintf("Executed command: %s", c.name))
			}
			break
		}
	}
}

func evalPerms(session disgord.Session, perms disgord.PermissionBits) (uint64, error) {
	guildID := disgord.ParseSnowflakeString(config.GuildID)
	user, err := session.GetCurrentUser(background)
	if err != nil {
		return 0, fmt.Errorf("Failed to get current user for permissions calculation; %w", err)
	}
	p, err := session.GetMemberPermissions(background, guildID, user.ID)
	if err != nil {
		return 0, fmt.Errorf("Failed to get permissions for current user; %w", err)
	}
	return p & perms, nil
}

func cleanup(session disgord.Session, duration int, ms ...*disgord.Message) {
	<-time.NewTimer(time.Second * time.Duration(duration)).C
	var p uint64
	for _, m := range ms {
		if m.IsDirectMessage() {
			continue
		}
		if p == 0 {
			p, err := evalPerms(session, disgord.PermissionManageMessages)
			if err != nil {
				logger.Error("Failed to delete message; %w", err)
				break
			} else if p == 0 {
				logger.Error("Failed to delete message: bot lacks permissions; %w", err)
				break
			}
		}
		if err := session.DeleteMessage(background, m.ChannelID, m.ID); err != nil {
			logger.Error("Failed to delete message; %w")
		}

	}

}
