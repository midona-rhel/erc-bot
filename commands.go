package main

import "github.com/bwmarrin/discordgo"

import "strings"

import "time"

func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iamnot") {
		b.removeRole(m, s)
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iam") {
		b.addRole(m, s)
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"help") {
		b.help(m, s)
	}
}

func (b *Bot) removeRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				if !userHasRole(r.RoleID, m.Member) {
					b.respondAndDelete("You do not have the role", m.ChannelID, time.Second*30)
					return
				}
				err := s.GuildMemberRoleRemove(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					logRemoveRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					logCommand(m, "removeRole")
					b.respondAndDelete("Role removed", m.ChannelID, time.Second*30)
				}
			}
		}
	}
}

func (b *Bot) addRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				if userHasRole(r.RoleID, m.Member) {
					b.respondAndDelete("You already have the role", m.ChannelID, time.Second*30)
					return
				}
				err := s.GuildMemberRoleAdd(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					logAddRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					logCommand(m, "addRole")
					b.respondAndDelete("Role added", m.ChannelID, time.Second*30)
				}
			}
		}
	}
}

func (b *Bot) help(m *discordgo.MessageCreate, s *discordgo.Session) error {
	logCommand(m, "help")
	_, err := s.ChannelMessage(m.ChannelID, b.config.Help)
	return err
}

func userHasRole(roleID string, member *discordgo.Member) bool {
	for _, r := range member.Roles {
		if r == roleID {
			return true
		}
	}
	return false
}
