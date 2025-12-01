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

//go:embed database/sqlite/migrations/*.sql
var sqliteMigrations embed.FS

//go:embed database/postgres/migrations/*.sql
var postgresMigrations embed.FS

//go:embed assets/*
var assetsFS embed.FS

func main() {
	config.Init()

	logger := config.InitLogger()

	database, store := config.InitDB(config.Migrations{
		config.SqliteDriver:   sqliteMigrations,
		config.PostgresDriver: postgresMigrations,
	})

	handler := handlers.New(&handlers.HandlerParams{
		Production:     config.Environment[config.EnvProd] == "1",
		Logger:         logger,
		Database:       database,
		TranslatorFunc: config.InitTranslator(),
		Sessions:       config.InitSessions(store),
		Secret:         config.Environment[config.EnvSecret],
	})

	routes := router.Init(handler, assetsFS)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info(
			"server starting",
			"address", ":"+config.Environment[config.EnvPort],
		)

		if err := http.ListenAndServe(
			":"+config.Environment[config.EnvPort],
			routes,
		); err != nil {
			logger.Error(
				"failed to start server",
				"port", config.Environment[config.EnvPort],
				"error", err,
			)

			os.Exit(1)
		}
	}()

	defer database.Close(context.Background())

	<-stop

	logger.Info("shutting down...")
}
