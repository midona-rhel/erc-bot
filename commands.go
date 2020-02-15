package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}
	if b.validCommand("iamnot", m) {
		b.removeRole(m, s)
	} else if b.validCommand("iam", m) {
		b.addRole(m, s)
	} else if b.validCommand("help", m) {
		b.help(m, s)
	} else if b.validCommand("check", m) {
		b.check(m, s)
	}
}

func (b *Bot) validCommand(command string, m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, b.config.CommandPrefix+command)
}

func (b *Bot) removeRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	message := strings.ToLower(m.Content)
	message = strings.Replace(message, b.config.CommandPrefix+"iam", "", 1)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				member, err := b.session.State.Member(b.config.Discord.DefaultGuild, m.Author.ID)
				if err != nil {
					log.Error(err)
					return
				}

				if !userHasRole(r.RoleID, member.Roles) {
					b.replyAndClear(fmt.Sprintf("You are not in %s", r.Alias[0]), m.ChannelID, m.ID, time.Second*30)
					return
				}
				err = s.GuildMemberRoleRemove(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					b.logRemoveRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					b.logCommand(m, fmt.Sprintf("Removed role %s", r.Alias[0]))
					b.replyAndClear(fmt.Sprintf("You have been removed from %s", r.Alias[0]), m.ChannelID, m.ID, time.Second*30)
					return
				}
			}
		}
	}
	b.replyAndClear(fmt.Sprintf("Sorry, the role %s does not exist", strings.TrimSpace(message)), m.ChannelID, m.ID, time.Second*30)
}

func (b *Bot) addRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	message := strings.ToLower(m.Content)
	message = strings.Replace(message, b.config.CommandPrefix+"iam", "", 1)
	for _, r := range b.config.Role {
		for _, alias := range r.Alias {
			if strings.Contains(message, alias) {
				member, err := b.session.State.Member(b.config.Discord.DefaultGuild, m.Author.ID)
				if err != nil {
					log.Error(err)
					return
				}
				if userHasRole(r.RoleID, member.Roles) {
					b.replyAndClear(fmt.Sprintf("You are already in %s", r.Alias[0]), m.ChannelID, m.ID, time.Second*30)
					return
				}
				err = s.GuildMemberRoleAdd(b.config.Discord.DefaultGuild, m.Author.ID, r.RoleID)
				if err != nil {
					b.logAddRoleError(m.Author.ID, b.config.Discord.DefaultGuild, r.RoleID, err)
				} else {
					b.logCommand(m, fmt.Sprintf("Added role %s", r.Alias[0]))
					b.replyAndClear(fmt.Sprintf("You now have the %s role", r.Alias[0]), m.ChannelID, m.ID, time.Second*30)
					return
				}
			}
		}
	}
	b.replyAndClear(fmt.Sprintf("Sorry, the role %s does not exist", message), m.ChannelID, m.ID, time.Second*30)
}

func (b *Bot) check(m *discordgo.MessageCreate, s *discordgo.Session) {

	if m.GuildID != "" {
		return
	}
	message, err := m.ContentWithMoreMentionsReplaced(s)
	message = strings.ToLower(message)
	message = strings.TrimSpace(strings.Replace(message, b.config.CommandPrefix+"check", "", 1))
	if err != nil {
		log.Error(err)
		return
	}
	regex := regexp.MustCompile("\\<(.*?)\\>")
	message = string(regex.ReplaceAll([]byte(message), []byte("")))

	newlines := strings.Count(message, "\n")
	characters := len(message)

	reply := fmt.Sprintf("Your message has %d characters and %d newlines\n", characters, newlines)
	for _, c := range b.config.Throttle {
		ch, err := s.Channel(c.ChannelID)
		if err != nil {
			log.Error(fmt.Errorf("Failed to get channel; %w", err))
			return
		}

		if (c.CharLimit <= characters && c.CharLimit > 0) || (c.NewlineLimit <= newlines && c.NewlineLimit > 0) {
			reply = reply + constructUnsuccesfulCheckRespose(ch.Name, c.CharLimit, c.NewlineLimit)
		} else {
			reply = reply + constructSuccesfulCheckRespose(ch.Name, c.CharLimit, c.NewlineLimit)
		}
	}
	b.pmUser(m.Author.ID, reply)
}

func constructUnsuccesfulCheckRespose(channel string, characters, newlines int) string {
	reply := fmt.Sprintf("❌Your message does not meet the requirement for channel %s, ", channel)
	if characters <= 0 {
		reply = reply + "there is no character limit"
	} else {
		reply = reply + fmt.Sprintf("the character limit is %d", characters)
	}

	if newlines <= 0 {
		reply = reply + " and there is no newline limit\n"
	} else {
		reply = reply + fmt.Sprintf(" and the newline limit is %d\n", newlines)
	}

	return reply
}

func constructSuccesfulCheckRespose(channel string, characters, newlines int) string {
	reply := fmt.Sprintf("✅Your message does meet the requirement for channel %s, ", channel)
	if characters <= 0 {
		reply = reply + "there is no character limit"
	} else {
		reply = reply + fmt.Sprintf("the character limit is %d", characters)
	}

	if newlines <= 0 {
		reply = reply + " and there is no newline limit\n"
	} else {
		reply = reply + fmt.Sprintf(" and the newline limit is %d\n", newlines)
	}

	return reply
}

func (b *Bot) help(m *discordgo.MessageCreate, s *discordgo.Session) {
	b.logCommand(m, "help")
	b.replyAndClear(b.config.Help, m.ChannelID, m.ID, time.Second*30)
}

func userHasRole(roleID string, roles []string) bool {
	for _, r := range roles {
		if r == roleID {
			return true
		}
	}
	return false
}
