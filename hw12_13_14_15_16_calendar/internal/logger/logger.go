package logger

import (
	"os"

	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/sirupsen/logrus" //nolint: depguard
)

type logger struct {
	level string
}

func New(level string) i.Logger {
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

	return &logger{
		level: level,
	}
}

func (logger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (logger) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (logger) Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func (logger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (logger) Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
