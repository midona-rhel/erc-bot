package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	monitor = logrus.New()
)

func initMonitor(filename string) {
	f, err := os.OpenFile("./"+filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	monitor.SetFormatter(&logrus.JSONFormatter{})
	monitor.SetOutput(f)
}

func monitorGuildAdd(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member join")
}

func monitorGuildRemove(_ *discordgo.Session, g *discordgo.GuildMemberRemove) {
	monitor.WithFields(logrus.Fields{
		"userID":  g.Member.User.ID,
		"guildID": g.Member.GuildID,
	}).Info("guild member remove")
}

func monitorMessageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	} else {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}
}

func monitorMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	if m.Author == nil {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    "bot",
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	} else {
		monitor.WithFields(logrus.Fields{
			"content":   m.Content,
			"userID":    m.Author.ID,
			"channelID": m.ChannelID,
			"guildID":   m.GuildID,
		}).Info("message delete")
	}
}

func monitorMessageUpdate(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author == nil {
		monitor.WithFields(logrus.Fields{
			"content":          m.Content,
			"previous content": m.BeforeUpdate.Content,
			"userID":           "bot",
			"channelID":        m.ChannelID,
			"guildID":          m.GuildID,
		}).Info("message delete")
	} else {
		monitor.WithFields(logrus.Fields{
			"content":          m.Content,
			"previous content": m.BeforeUpdate.Content,
			"userID":           m.Author.ID,
			"channelID":        m.ChannelID,
			"guildID":          m.GuildID,
		}).Info("message delete")
	}
}
