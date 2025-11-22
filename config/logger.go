package config

import (
	"log/slog"
	"os"
	"strconv"
)

func InitLogger() *slog.Logger {
	production := Environment[EnvProd] == "1"
	levelStr := Environment[EnvLog]

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		level = int(slog.LevelInfo)
	}

	logOpts := &slog.HandlerOptions{
		AddSource: production,
		Level:     slog.Level(level),
	}

	if production {
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
