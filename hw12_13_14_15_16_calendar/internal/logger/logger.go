package logger

import (
	"os"

	"github.com/sirupsen/logrus" //nolint: depguard
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetOutput(os.Stdout)

	return &Logger{
		level: level,
	}
}

func (Logger) Debug(msg string) {
	logrus.Debug(msg)
}

func (Logger) Info(msg string) {
	logrus.Info(msg)
}

func (Logger) Warn(msg string) {
	logrus.Warn(msg)
}

func (Logger) Error(msg string) {
	logrus.Error(msg)
}
