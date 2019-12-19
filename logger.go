package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	log   = logrus.New()
	fmter = new(logrus.TextFormatter)
)

func initLog() {
	f, err := os.OpenFile("./log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(f)
	log.AddHook(new(Hook))
}

// Hook is a logrus hook.
type Hook struct{}

// Levels return which levels the hook fires on.
func (h *Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.InfoLevel,
	}
}

// Fire is the action taken when a correct level event happens.
func (h *Hook) Fire(entry *logrus.Entry) (err error) {
	line, err := fmter.Format(entry)
	if err == nil {
		fmt.Fprintf(os.Stderr, string(line))
	}
	return
}

func logCommand(m *discordgo.MessageCreate, command string) {
	log.WithFields(logrus.Fields{
		"command":    command,
		"userInput":  m.Content,
		"userID":     m.Author.ID,
		"channeldID": m.ChannelID,
	}).Info("executed command")
}

func logMessageSendError(channelID string, err error) {
	log.WithFields(logrus.Fields{
		"channelID": channelID,
	}).Error(err)
}

func logMessageDeleteError(channelID, messageID string, err error) {
	log.WithFields(logrus.Fields{
		"channelID": channelID,
		"messageID": messageID,
	}).Error(err)
}

func logMessagePurging(amount int, channelID string) {
	log.WithFields(logrus.Fields{
		"amount":    amount,
		"channelID": channelID,
	}).Info("purged channel")
}
func logMessagePurgingError(channelID string, err error) {
	log.WithFields(logrus.Fields{
		"action":    "purge channel",
		"channelID": channelID,
	}).Error(err)
}

func logRemoveRoleError(userID, guildID, roleID string, err error) {
	log.WithFields(logrus.Fields{
		"roleID":  roleID,
		"userID":  userID,
		"guildID": guildID,
	}).Error("failed to remove role:", err)
}

func logAddRoleError(userID, guildID, roleID string, err error) {
	log.WithFields(logrus.Fields{
		"roleID":  roleID,
		"userID":  userID,
		"guildID": guildID,
	}).Error("failed to add role:", err)
}
