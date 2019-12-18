package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"time"
)

func (b *Bot) purge(s *discordgo.Session) {
	cron := cron.New()
	for _, c := range b.config.Purge {
		cron.AddFunc(c.CronExpression, func() {
			ms, err := s.ChannelMessages(c.ChannelID, 100, "", "", "")
			if err != nil {
				fmt.Println(err)
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
				fmt.Println(err)
			}
		})
	}
}
