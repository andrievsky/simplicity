package loggers

import (
	"log/slog"
	"os"
	"simplicity/config"
)

func NewLogger(conf *config.Config) *slog.Logger {
	level := slog.LevelInfo
	var ra = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{} // remove time
		}
		return a
	}
	if conf.EnableDebug {
		level = slog.LevelDebug
		ra = nil
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   conf.EnableDebug,
		Level:       level,
		ReplaceAttr: ra,
	})

	return slog.New(handler)
}
