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

func (b *Bot) logCommand(m *discordgo.MessageCreate, command string) {
	log.WithFields(logrus.Fields{
		"command":    command,
		"userInput":  m.Content,
		"userID":     m.Author.ID,
		"channeldID": m.ChannelID,
	}).Info("executed command")
	b.session.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("Executed command %s on user %s", command, m.Author.Username+"#"+m.Author.ID))
}

func (b *Bot) logMessageSendError(channelID string, err error) {
	log.WithFields(logrus.Fields{
		"channelID": channelID,
	}).Error(err)
}

func (b *Bot) logMessageDeleteError(channelID, messageID string, err error) {
	log.WithFields(logrus.Fields{
		"channelID": channelID,
		"messageID": messageID,
	}).Error(err)
}

func (b *Bot) logMessagePurging(amount int, channelID string) {
	log.WithFields(logrus.Fields{
		"amount":    amount,
		"channelID": channelID,
	}).Info("purged channel")
	c, err := b.session.Channel(channelID)
	if err != nil {
		log.Error(err)
		return
	}
	b.session.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("Purged %d messages in channel %s", amount, c.Name+":"+channelID))
}
func (b *Bot) logMessagePurgingError(channelID string, err error) {
	log.WithFields(logrus.Fields{
		"action":    "purge channel",
		"channelID": channelID,
	}).Error(err)
}

func (b *Bot) logRemoveRoleError(userID, guildID, roleID string, err error) {
	log.WithFields(logrus.Fields{
		"roleID":  roleID,
		"userID":  userID,
		"guildID": guildID,
	}).Error("failed to remove role:", err)
}

func (b *Bot) logAddRoleError(userID, guildID, roleID string, err error) {
	log.WithFields(logrus.Fields{
		"roleID":  roleID,
		"userID":  userID,
		"guildID": guildID,
	}).Error("failed to add role:", err)
}

func (b *Bot) logThrottleUser(m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}
	log.WithFields(logrus.Fields{
		"userID":    m.Author.ID,
		"channelID": m.ChannelID,
	}).Info("throttled user")
	c, err := b.session.Channel(m.ChannelID)
	if err != nil {
		log.Error(err)
		return
	}
	b.session.ChannelMessageSend(b.config.Monitor.Output, fmt.Sprintf("Throttled user %s in channel %s", m.Author.Username+"#"+m.Author.ID, c.Name+"."+m.ChannelID))
}

func (b *Bot) logFailedToCreateChannel(userID string, err error) {
	log.WithFields(logrus.Fields{
		userID: userID,
	}).Error("failed to create private channel:", err)
}
