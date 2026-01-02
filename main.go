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

	h := handler.New(&handler.HandlerParams{
		Logger: logger,
	})

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
