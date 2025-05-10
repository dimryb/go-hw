package logger

import "fmt"

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{
		level: level,
	}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
