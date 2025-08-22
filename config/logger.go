package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

const (
	defaultSeverity = 0
	// LevelDebug = -4
	// LevelInfo  = 0
	// LevelWarn  = 4
	// LevelError = 8
)

func InitLogger(production bool, severityStr string) (*slog.Logger, error) {
	severity := defaultSeverity
	if severityStr != "" {
		var err error
		severity, err = strconv.Atoi(severityStr)
		if err != nil {
			return nil, fmt.Errorf("error converting severity string to int: %v", err)
		}
	}

	logOpts := &slog.HandlerOptions{
		AddSource: !production,
		Level:     slog.Level(severity),
	}

	var logger *slog.Logger
	if production {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, logOpts))
	}

	return logger, nil
}
