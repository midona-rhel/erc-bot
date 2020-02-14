package main

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var (
	config     = loadConfig()
	background = context.Background()
	logger     = &Logger{
		debug:        true,
		logChannelID: disgord.ParseSnowflakeString("656571472705224704"),
		logger:       logrus.New(),
	}
)

func main() {
	session := disgord.New(disgord.Config{
		LoadMembersQuietly: true,
		BotToken:           config.Token,
		Logger:             logger,
	})
	logger.session = session
	session.On(disgord.EvtReady, func() {
		logger.ready = true
		logger.Debug("Hello World")
	})
	session.On(disgord.EvtMessageCreate, commandHandler)
	session.StayConnectedUntilInterrupted(background)
}
