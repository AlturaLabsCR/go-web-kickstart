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

//go:embed assets/*
var assetsFS embed.FS

func main() {
	config.Init()

	logger := config.InitLogger()
	database := config.InitDB()
	tr := config.InitTranslator()

	handler := handlers.New(&handlers.HandlerParams{
		Production:     config.Environment[config.EnvProd] == "1",
		Logger:         logger,
		Database:       database,
		TranslatorFunc: tr,
	})

	routes := router.Init(handler, assetsFS)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", "address", ":"+config.Environment[config.EnvPort])

		if err := http.ListenAndServe(":"+config.Environment[config.EnvPort], routes); err != nil {
			logger.Error("failed to start server", "port", config.Environment[config.EnvPort], "error", err)
			os.Exit(1)
		}
	}()

	defer database.Close(context.Background())

	<-stop

	logger.Info("shutting down...")
}
