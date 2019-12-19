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
		"username": m.Author.Username,
		"nickname": m.Member.Nick,
		"userID":   m.Author.ID,
		"guildID":  m.GuildID,
	}).Info("guild member remove")
}
