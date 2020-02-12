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

// Logger is the custom logger for the erc-bot
type Logger struct {
	logger       *logrus.Logger
	session      disgord.Session
	logChannelID disgord.Snowflake
	debug        bool
	ready        bool
}

// Debug only logs if debug is enabled
func (l *Logger) Debug(v ...interface{}) {
	if l.debug {
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
}

// Error logs errors
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

// Info logs info, like other methods for logger messages are sent to a discord channel when the logger is active.
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

// ErrorQuiet logs the error without it being sent to a discord channel, this is called from Error if it fails to send a
// message to said discord channel
func (l *Logger) ErrorQuiet(v ...interface{}) {
	l.Error(v...)
}
