package logger

import (
	"log"
	"log/slog"
	"os"
)

type logger struct {
	logger *slog.Logger
}

var _ Logger = (*logger)(nil)

func New(level string) Logger {
	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseSlogLevel(level),
	})

	return logger{
		logger: slog.New(loggerHandler),
	}
}

func parseSlogLevel(level string) slog.Leveler {
	switch level {
	case "error":
		return slog.LevelError
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	}

	return slog.LevelDebug
}

func (l logger) Named(name string) Logger {
	handler := l.logger.Handler().WithAttrs([]slog.Attr{
		{
			Key:   "name",
			Value: slog.StringValue(name),
		},
	})

	return logger{
		logger: slog.New(handler),
	}
}

func (l logger) Debug(message string, args ...interface{}) {
	l.logger.Debug(message, args...)
}

func (l logger) Info(message string, args ...interface{}) {
	l.logger.Info(message, args...)
}

func (l logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, args...)
}

func (l logger) Error(message string, args ...interface{}) {
	l.logger.Error(message, args...)
}

func (l logger) Fatal(message string, args ...interface{}) {
	log.Fatal(message, args)
}
