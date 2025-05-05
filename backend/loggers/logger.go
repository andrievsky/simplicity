package loggers

import (
	"log/slog"
	"os"
	"simplicity/config"
)

func NewLogger(conf *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if conf.EnableDebug {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	return slog.New(handler)
}
