package main

import "github.com/bwmarrin/discordgo"

import "strings"

func (b *Bot) handleRoles(m *discordgo.MessageCreate, s *discordgo.Session) {
	if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iamnot") {
		b.removeRole(m, s)
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iam") {
		b.addRole(m, s)
	}
}

func (b *Bot) removeRole(m *discordgo.MessageCreate, s *discordgo.Session) error {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				return s.GuildMemberRoleRemove(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
			}
		}
	}
	return nil
}

func (b *Bot) addRole(m *discordgo.MessageCreate, s *discordgo.Session) error {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				return s.GuildMemberRoleAdd(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
			}
		}
	}
	return nil
}
