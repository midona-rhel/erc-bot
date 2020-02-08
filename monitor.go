package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	monitor = logrus.New()
)

func initMonitor() {
	f, err := os.OpenFile("./chatlog.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	monitor.SetFormatter(&logrus.JSONFormatter{})
	monitor.SetOutput(f)
}

func (b *Bot) monitorGuildAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	t := time.Now()
	title := fmt.Sprintf("[%d:%d:%d]**New User**", t.Hour(), t.Minute(), t.Second())
	b.sengLogMessage(title, getUserName(g.User, g.Member), g.User.AvatarURL(""), "")
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member join")
}

func (b *Bot) monitorGuildRemove(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	t := time.Now()
	title := fmt.Sprintf("[%d:%d:%d]**New Left**", t.Hour(), t.Minute(), t.Second())
	b.sengLogMessage(title, getUserName(g.User, g.Member), g.User.AvatarURL(""), "")
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member remove")
}

func (b *Bot) monitorMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message create")
	} else {
		t := time.Now()
		name := getUserName(m.Author, m.Member)
		content := fmt.Sprintf("[%d:%d:%d]**New Message** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
		description, _ := m.ContentWithMoreMentionsReplaced(s)
		b.sengLogMessage(name, content, m.Author.AvatarURL("50x50"), description)
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message create")
	}
}

func (b *Bot) monitorMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	a, err := b.session.GuildAuditLog(m.GuildID, "", "", discordgo.AuditLogActionMessageDelete, 1)
	if err != nil || len(a.AuditLogEntries) != 1 {
		return
	}
	var audit = a.AuditLogEntries[0]
	t := time.Now()
	user, err := s.User(audit.UserID)
	if err != nil {
		return
	}
	content := fmt.Sprintf("[%d:%d:%d]**Message Deleted** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
	b.sengLogMessage(user.Username, content, user.AvatarURL(""), "Unkown content")
	monitor.WithFields(logrus.Fields{
		"content":   m.Content,
		"channelID": m.ChannelID,
		"guildID":   m.GuildID,
	}).Info("message delete")
}

func (b *Bot) monitorMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author == nil || m.Author.Bot {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	} else {
		t := time.Now()
		name := getUserName(m.Author, m.Member)
		content := fmt.Sprintf("[%d:%d:%d]**Message Updated** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
		description, _ := m.ContentWithMoreMentionsReplaced(s)
		b.sengLogMessage(name, content, m.Author.AvatarURL(""), description)
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}

}

func (b *Bot) sengLogMessage(author, content, icon, field string) {
	b.session.ChannelMessageSendComplex(b.config.Monitor.Output, &discordgo.MessageSend{
		Content: content,
		Embed: &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: icon,
				Name:    author,
			},
			Description: field,
		},
	})
}

func getUserName(u *discordgo.User, m *discordgo.Member) string {
	if m == nil {
		if u == nil {
			return ""
		}
		return u.Username
	} else if m.Nick != "" {
		return m.Nick
	}
	return u.Username
}

func (b *Bot) getChannelName(channelID string) string {
	c, err := b.session.Channel(channelID)
	if err != nil {
		log.Error(err)
		return ""
	}
	if c.Name == "" {
		return "PM"
	}
	return c.Mention()
}
