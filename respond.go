package main

import "github.com/bwmarrin/discordgo"

import "time"

func (b *Bot) respond(content, channelID string) (*discordgo.Message, error) {
	m, err := b.session.ChannelMessageSend(channelID, content)
	if err != nil {
		logMessageSendError(channelID, err)
	}
	return m, err
}

func (b *Bot) respondAndDelete(content, channelID, messageID string, t time.Duration) {
	m, err := b.respond(content, channelID)
	if err != nil {
		return
	}
	time.AfterFunc(t, func() {
		err := b.session.ChannelMessageDelete(channelID, m.ID)
		if err != nil {
			logMessageDeleteError(channelID, m.ID, err)
		}
		err = b.session.ChannelMessageDelete(channelID, messageID)
		if err != nil {
			logMessageDeleteError(channelID, m.ID, err)
		}
	})
}
