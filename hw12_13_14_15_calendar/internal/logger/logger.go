package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string, output string) *Logger {
	var logger *slog.Logger
	var l slog.Level

	err := l.UnmarshalText([]byte(level))
	if err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	options := &slog.HandlerOptions{
		Level: l,
	}
	switch output {
	case "stdout", "":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, options))

	default:
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			log.Fatalf("Failed to create/open a log file: %v", err)
		}
		logger = slog.New(slog.NewJSONHandler(file, options))
	}

	return &Logger{
		logger: logger,
	}
}

func (l Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l Logger) Debug(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
