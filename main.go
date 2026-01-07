package main

import (
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/config"
	"app/handler"
	"app/router"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {
	config.InitEnv()
	config.InitRoutes()

	logger := config.InitLogger()
	logger.Debug("config", "config", config.Config)

	database, err := config.InitDB()
	if err != nil {
		logger.Error("failed to init database", "error", err)
		os.Exit(1)
	}

	storage, err := config.InitStorage()
	if err != nil {
		logger.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	translator, err := config.InitTranslator()
	if err != nil {
		logger.Error("failed to init translator", "error", err)
		os.Exit(1)
	}

	sessions, err := config.InitSessions(database)
	if err != nil {
		logger.Error("failed to init sessions", "error", err)
		os.Exit(1)
	}

	h, err := handler.New(&handler.HandlerParams{
		Logger:     logger,
		Database:   database,
		Storage:    storage,
		Translator: translator,
		Sessions:   sessions,
	})
	if err != nil {
		logger.Error("failed to init handler", "error", err)
		os.Exit(1)
	}

	routes := router.Init(h, assetsFS)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := config.Config.App.Port

		logger.Info("server starting", "address", ":"+port)

		if err := http.ListenAndServe(":"+port, routes); err != nil {
			logger.Error("failed to start server", "port", port, "error", err)
			os.Exit(1)
		}
	}()
	// defer database.Close(context.Background())

	<-stop
	logger.Info("shutting down...")
}
