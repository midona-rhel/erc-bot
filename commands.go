package main

import "github.com/bwmarrin/discordgo"

import "strings"

import "time"

func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if b.validCommand("iamnot", m) {
		b.removeRole(m, s)
	} else if b.validCommand("iam", m) {
		b.addRole(m, s)
	} else if b.validCommand("help", m) {
		b.help(m, s)
	}
}

func (b *Bot) validCommand(command string, m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, b.config.CommandPrefix + command)
}

func (b *Bot) removeRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				member, err := b.session.State.Member(b.config.Discord.DefaultGuild, m.Author.ID)
				if err != nil {
					log.Error(err)
					return
				}

				if !userHasRole(r.RoleID, member.Roles) {
					b.respondAndDelete("You do not have the role", m.ChannelID, m.ID, time.Second*30)
					return
				}
				err = s.GuildMemberRoleRemove(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					logRemoveRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					logCommand(m, "removeRole")
					b.respondAndDelete("Role removed", m.ChannelID, m.ID, time.Second*30)
					return
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
				member, err := b.session.State.Member(b.config.Discord.DefaultGuild, m.Author.ID)
				if err != nil {
					log.Error(err)
					return
				}
				if userHasRole(r.RoleID, member.Roles) {
					b.respondAndDelete("You already have the role", m.ChannelID, m.ID, time.Second*5)
					return
				}
				err = s.GuildMemberRoleAdd(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					logAddRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					logCommand(m, "addRole")
					b.respondAndDelete("Role added", m.ChannelID, m.ID, time.Second*5)
					return
				}
			}
		}
	}
}

func (b *Bot) help(m *discordgo.MessageCreate, s *discordgo.Session) {
	logCommand(m, "help")
	b.respondAndDelete(b.config.Help, m.ChannelID, m.ID, time.Hour)
}

func userHasRole(roleID string, roles []string) bool {
	for _, r := range roles {
		if r == roleID {
			return true
		}
	}
	return false
}
