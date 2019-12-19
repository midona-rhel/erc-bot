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
	monitor.SetOutput(f)
}

func monitorGuildAdd(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {
	monitor.WithFields(logrus.Fields{
		"username": g.User.Username,
		"userID":   g.Member.User.ID,
		"guildID":  g.Member.GuildID,
	}).Info("guild member join")
}

func monitorGuildRemove(_ *discordgo.Session, g *discordgo.GuildMemberRemove) {
	monitor.WithFields(logrus.Fields{
		"username": g.User.Username,
		"userID":   g.Member.User.ID,
		"guildID":  g.Member.GuildID,
	}).Info("guild member remove")
}

func monitorMessageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	monitor.WithFields(logrus.Fields{
		"content":   m.Content,
		"username":  m.Author.Username,
		"userID":    m.Author.ID,
		"channelID": m.ChannelID,
		"guildID":   m.GuildID,
		"bot":       m.Author.Bot,
	}).Info("message create")
}

func monitorMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	monitor.WithFields(logrus.Fields{
		"content":   m.Content,
		"username":  m.Author.Username,
		"userID":    m.Author.ID,
		"channelID": m.ChannelID,
		"guildID":   m.GuildID,
	}).Info("message delete")
}

func monitorMessageUpdate(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	monitor.WithFields(logrus.Fields{
		"content":          m.Content,
		"previous content": m.BeforeUpdate.Content,
		"username":         m.Author.Username,
		"userID":           m.Author.ID,
		"channelID":        m.ChannelID,
		"guildID":          m.GuildID,
	}).Info("message update")
}
