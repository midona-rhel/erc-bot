package main

import (
	"github.com/sirupsen/logrus"
)

// Logger represents a struct
type Logger struct {
	loggers []logrus.Logger
}

// Info logs info logs
func (l *Logger) Info(data ...interface{}) {

}

// Error logs error logs
func (l *Logger) Error(data ...interface{}) {
}
