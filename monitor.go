package main

import (
	"fmt"
	"os"

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
	s.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("User %s joined", g.Member.User.Username+"#"+g.Member.User.ID))
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member join")
}

func (b *Bot) monitorGuildRemove(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	s.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("User %s left", g.Member.User.Username+"#"+g.Member.User.ID))
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
		s.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("User %s created message with content: %s", m.Author.Username+"#"+m.Author.ID, m.ContentWithMentionsReplaced()))
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message create")
	}
}

func (b *Bot) monitorMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.Author == nil || m.Author.Bot {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	} else {
		s.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("User %s deleted message with content: %s", m.Author.Username+"#"+m.Author.ID, m.ContentWithMentionsReplaced()))
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}
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
		s.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("User %s deleted message with content: %s", m.Author.Username+"#"+m.Author.ID, m.ContentWithMentionsReplaced()))
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}
}
