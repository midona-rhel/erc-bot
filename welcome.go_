package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleWelcomeMessage(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {
	s := strings.Replace(b.config.WelcomeMessage, ":NAME:", g.User.Username, 1)
	b.pmUser(g.User.ID, s)
}
