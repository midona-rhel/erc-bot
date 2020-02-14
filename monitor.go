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
	content := fmt.Sprintf("[%02d:%02d:%02d] **User Joined**", t.Hour(), t.Minute(), t.Second())
	b.sendLogMessage(getUserName(g.User, g.Member), content, g.User.AvatarURL(""), "", 0)
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member join")
}

func (b *Bot) monitorGuildRemove(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	t := time.Now()
	content := fmt.Sprintf("[%02d:%02d:%02d] **User Left**", t.Hour(), t.Minute(), t.Second())
	b.sendLogMessage(getUserName(g.User, g.Member), content, g.User.AvatarURL(""), "", 0)
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member remove")
}

func (b *Bot) monitorMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.messages.setMessage(m.ID, m)
	if m.Author == nil || m.Author.Bot {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message create")
	} else {
		color := 0
		if len(m.MentionRoles)+len(m.Mentions) > 0 {
			color = 0x22AA22
		}
		t := time.Now()
		name := getUserName(m.Author, m.Member)
		content := fmt.Sprintf("[%02d:%02d:%02d] **New Message** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
		description, _ := m.ContentWithMoreMentionsReplaced(s)
		b.sendLogMessage(name, content, m.Author.AvatarURL(""), description, color)
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message create")
	}
}

func (b *Bot) monitorMessageDelete(s *discordgo.Session, d *discordgo.MessageDelete) {
	m, ok := b.messages.getMessage(d.Message.ID)
	if !ok {
		return
	}
	t := time.Now()
	name := getUserName(m.Author, m.Member)
	content := fmt.Sprintf("[%02d:%02d:%02d ] **Message Deleted** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
	description, _ := m.ContentWithMoreMentionsReplaced(s)
	b.sendLogMessage(name, content, m.Author.AvatarURL(""), description, 0)
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
		content := fmt.Sprintf("[%02d:%02d:%02d] **Message Updated** %s", t.Hour(), t.Minute(), t.Second(), b.getChannelName(m.ChannelID))
		description, _ := m.ContentWithMoreMentionsReplaced(s)
		b.sendLogMessage(name, content, m.Author.AvatarURL(""), description, 0)
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}

}

func (b *Bot) sendLogMessage(author, content, icon, field string, color int) {
	b.session.ChannelMessageSendComplex(b.config.Monitor.Output, &discordgo.MessageSend{
		Content: content,
		Embed: &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: icon,
				Name:    author,
			},
			Description: field,
			Color:       color,
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
