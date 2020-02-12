package main

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var (
	colorRed    = 0xFF0000
	colorGreen  = 0x00FF00
	colorYellow = 0x888800
	colorBlue   = 0x0000FF
)

// Logger is the custom logger for erc-bot
type Logger struct {
	logger       *logrus.Logger
	session      disgord.Session
	logChannelID disgord.Snowflake
	debug        bool
	ready        bool
}

func (l *Logger) Debug(v ...interface{}) {
	if l.ready {
		_, err := l.session.SendMsg(context.Background(), l.logChannelID,
			disgord.Embed{
				Title:       "Debug Message",
				Description: fmt.Sprint(v...),
				Color:       colorYellow,
			})
		if err != nil {
			l.Error(err)
		}
	}
	l.logger.Debug(v...)
}

func (l *Logger) Error(v ...interface{}) {
	if l.ready {
		_, err := l.session.SendMsg(context.Background(), l.logChannelID,
			disgord.Embed{
				Title:       "Error Message",
				Description: fmt.Sprint(v...),
				Color:       colorRed,
			})
		if err != nil {
			l.ErrorQuiet(err)
		}
	}
	l.logger.Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	if l.ready {

		_, err := l.session.SendMsg(context.Background(), l.logChannelID,
			disgord.Embed{
				Title:       "Info Message",
				Description: fmt.Sprint(v...),
				Color:       colorBlue,
			})
		if err != nil {
			l.Error(err)
		}
	}
	l.logger.Debug(v...)
}

func (l *Logger) ErrorQuiet(v ...interface{}) {
	l.Error(v...)
}
