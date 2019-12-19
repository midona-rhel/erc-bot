package main

import "github.com/bwmarrin/discordgo"

import "strings"

func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error
	if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iamnot") {
		err = b.removeRole(m, s)
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"iam") {
		err = b.addRole(m, s)
	} else if strings.HasPrefix(m.Content, b.config.CommandPrefix+"help") {
		err = b.help(m, s)
	}
	if err != nil {
		log.Error(err)
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
