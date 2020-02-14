package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func (b *Bot) purge(s *discordgo.Session) {
	for _, c := range b.config.Purge {
		cr := cron.New()
		channelID := c.ChannelID
		cr.AddFunc(c.CronExpression, b.cronHandler(channelID))
		cr.Start()
	}

}

func (b *Bot) cronHandler(channelID string) func() {
	return func() {
		date := time.Now().AddDate(0, 0, -13)
		ms, err := b.session.ChannelMessages(channelID, 100, "", "", "")
		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "get channel messages",
				"channelID": channelID,
			}).Error(err)
		}

		var messagesToDelete []string
		for _, m := range ms {
			t, err := m.Timestamp.Parse()
			if err != nil {
				panic(err)
			}
			if !t.Add(time.Second).After(time.Now()) && t.After(date) {
				messagesToDelete = append(messagesToDelete, m.ID)
			}
		}
		if len(messagesToDelete) == 0 {
			return
		}
		if err = b.session.ChannelMessagesBulkDelete(channelID, messagesToDelete); err != nil {
			b.logMessagePurgingError(channelID, err)
		} else {
			b.logMessagePurging(len(messagesToDelete), channelID)
		}
	}
}
