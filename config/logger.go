package config

import (
	"log/slog"
	"os"
	"strconv"
)

func InitLogger() *slog.Logger {
	levelStr := Config.App.LogLevel

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		level = int(slog.LevelInfo)
	}

	logOpts := &slog.HandlerOptions{
		AddSource: level < 0,
		Level:     slog.Level(level),
	}

	if level >= 0 {
		return slog.New(slog.NewJSONHandler(
			os.Stdout,
			logOpts,
		))
	}

	return slog.New(slog.NewTextHandler(
		os.Stdout,
		logOpts,
	))
}
