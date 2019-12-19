package main

import "github.com/bwmarrin/discordgo"

import "github.com/sirupsen/logrus"

func snoopGuildAdd(g *discordgo.GuildMemberAdd, _ *discordgo.Session) {
	log.WithFields(logrus.Fields{
		"username": g.User.Username,
		"userID":   g.Member.User.ID,
		"guildID":  g.Member.GuildID,
	}).Info("guild member join")
}

func snoopGuildRemove(g *discordgo.GuildMemberRemove, _ *discordgo.Session) {
	log.WithFields(logrus.Fields{
		"username": g.User.Username,
		"userID":   g.Member.User.ID,
		"guildID":  g.Member.GuildID,
	}).Info("guild member remove")
}

func snoopMessageCreate(m *discordgo.MessageCreate, _ *discordgo.Session) {
	log.WithFields(logrus.Fields{
		"username":  m.Author.Username,
		"userID":    m.Author.ID,
		"channelID": m.ChannelID,
		"guildID":   m.GuildID,
		"bot":       m.Author.Bot,
		"content":   m.Content,
	}).Info("message create")
}

func snoopMessageDelete(m *discordgo.MessageDelete, _ *discordgo.Session) {
	log.WithFields(logrus.Fields{
		"username":  m.Author.Username,
		"userID":    m.Author.ID,
		"channelID": m.ChannelID,
		"guildID":   m.GuildID,
		"content":   m.Content,
	}).Info("message delete")
}

func snoopMessageUpdate(m *discordgo.MessageUpdate, _ *discordgo.Session) {
	log.WithFields(logrus.Fields{
		"username":         m.Author.Username,
		"userID":           m.Author.ID,
		"channelID":        m.ChannelID,
		"guildID":          m.GuildID,
		"content":          m.Content,
		"previous content": m.BeforeUpdate.Content,
	}).Info("message update")
}
