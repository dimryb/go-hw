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

func (Logger) Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (Logger) Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (Logger) Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func (Logger) Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (Logger) Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
