package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/config"
	"app/handlers"
	"app/router"
)

//go:embed database/sqlite/migrations/*.up.sql
var sqliteMigrations embed.FS

//go:embed database/postgres/migrations/*.up.sql
var postgresMigrations embed.FS

//go:embed assets/*
var assetsFS embed.FS

func main() {
	config.Init()

	logger := config.InitLogger()

	logger.Debug("configuration", "config", config.Config)

	database, store := config.InitDB(config.Migrations{
		config.SqliteDriver:   sqliteMigrations,
		config.PostgresDriver: postgresMigrations,
	})

	handler := handlers.New(&handlers.HandlerParams{
		Logger:         logger,
		Database:       database,
		TranslatorFunc: config.InitTranslator(),
		Sessions:       config.InitSessions(store),
		Secret:         config.Config.App.Secret,
	})

	routes := router.Init(handler, assetsFS)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	port := config.Config.App.Port
	go func() {
		logger.Info(
			"server starting",
			"address", ":"+port,
		)

		if err := http.ListenAndServe(
			":"+port,
			routes,
		); err != nil {
			logger.Error(
				"failed to start server",
				"port", port,
				"error", err,
			)

			os.Exit(1)
		}
	}()

	defer database.Close(context.Background())

	<-stop

	logger.Info("shutting down...")
}
