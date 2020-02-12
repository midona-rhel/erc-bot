package main

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var (
	config     = loadConfig()
	background = context.Background()
)

func main() {
	session := disgord.New(disgord.Config{
		LoadMembersQuietly: true,
		BotToken:           config.Token,
		Logger:             new(logrus.Logger),
	})
	session.StayConnectedUntilInterrupted(background)
}
