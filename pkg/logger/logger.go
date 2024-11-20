package logger

import (
	"log/slog"
	"os"
)

type Log struct {
	original *slog.Logger
}

func NewLogger() *Log {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &Log{original: logger}
}

func (l *Log) With(args ...any) *Log {
	logger := l.original.With(args...)
	return &Log{logger}
}

func (l *Log) Info(msg string, args ...any) {
	l.original.Info(msg, args...)
}

func (l *Log) Error(msg string, err error) {
	l.original.Error(msg, slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
}

func (l *Log) WithError(err error, msg string, args ...any) error {
	passArgs := make([]any, len(args)+2)
	passArgs[0] = "error"
	passArgs[1] = err
	for i, arg := range args {
		passArgs[i+2] = arg
	}
	l.original.Error(msg, passArgs...)
	return err
}
