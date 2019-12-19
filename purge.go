package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func (b *Bot) purge(s *discordgo.Session) {
	cron := cron.New()
	for _, c := range b.config.Purge {
		cron.AddFunc(c.CronExpression, func() {
			ms, err := s.ChannelMessages(c.ChannelID, 100, "", "", "")
			if err != nil {
				log.WithFields(logrus.Fields{
					"action":    "get channel messages",
					"channelID": c.ChannelID,
				}).Error(err)
			}

			var messagesToDelete []string
			for _, m := range ms {
				t, err := m.Timestamp.Parse()
				if err != nil {
					panic(err)
				}
				if !t.Add(time.Hour * 24).After(time.Now()) {
					messagesToDelete = append(messagesToDelete, m.ID)
				}
			}
			if err = s.ChannelMessagesBulkDelete(c.ChannelID, messagesToDelete); err != nil {
				logMessagePurgingError(c.ChannelID, err)
			} else {
				logMessagePurging(len(messagesToDelete), c.ChannelID)
			}
		})
	}
}
