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
	defer f.Close()
	log.SetOutput(f)
	log.AddHook(new(Hook))
}

// Hook is a logrus hook.
type Hook struct{}

// Levels return which levels the hook fires on.
func (h *Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
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
		"channeldID": m.ChannelID,
		"userID":     m.Author.ID,
		"userInput":  m.Content,
		"command":    command,
	}).Info("Executed command")
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
