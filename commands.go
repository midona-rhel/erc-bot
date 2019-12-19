package main

import "github.com/bwmarrin/discordgo"

import "strings"

import "github.com/sirupsen/logrus"

func (b *Bot) handleRoles(m *discordgo.MessageCreate, s *discordgo.Session) {
	if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iamnot") {
		log.Error(b.removeRole(m, s))
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iam") {
		log.Error(b.addRole(m, s))
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"help") {
		log.Error(b.help(m, s))
	}
}

func (b *Bot) removeRole(m *discordgo.MessageCreate, s *discordgo.Session) error {
	message := strings.ToLower(m.Content)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				logCommand(m, "removeRole")
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
				logCommand(m, "addRole")
				return s.GuildMemberRoleAdd(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
			}
		}
	}
	return nil
}

func (b *Bot) help(m *discordgo.MessageCreate, s *discordgo.Session) error {
	logCommand(m, "help")
	_, err := s.ChannelMessage(m.ChannelID, b.config.Help)
	return err
}

func logCommand(m *discordgo.MessageCreate, command string) {
	log.WithFields(logrus.Fields{
		"channeldID": m.ChannelID,
		"userID":     m.Author.ID,
		"userInput":  m.Content,
		"command":    command,
	}).Info("Executed command")
}
