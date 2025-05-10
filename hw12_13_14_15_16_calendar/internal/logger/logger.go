package logger

import (
	"os"

	"github.com/sirupsen/logrus"
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

func (_ Logger) Debug(msg string) {
	logrus.Debug(msg)
}

func (_ Logger) Info(msg string) {
	logrus.Info(msg)
}

func (_ Logger) Warn(msg string) {
	logrus.Warn(msg)
}

func (_ Logger) Error(msg string) {
	logrus.Error(msg)
}
