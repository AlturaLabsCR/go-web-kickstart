package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/config"
	"app/handler"
	"app/router"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.InitEnv()
	config.InitRoutes()

	logger := config.InitLogger()
	logger.Debug("config", "config", config.Config)

	database, err := config.InitDB(ctx)
	if err != nil {
		logger.Error("failed to init database", "error", err)
		os.Exit(1)
	}

	storage, err := config.InitStorage(ctx, database.Querier())
	if err != nil {
		logger.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	translator, err := config.InitTranslator()
	if err != nil {
		logger.Error("failed to init translator", "error", err)
		os.Exit(1)
	}

	sessions, err := config.InitSessions(ctx, database)
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

	port := config.Config.App.Port
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routes,
	}

	go func() {
		logger.Info("server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-stop

	logger.Debug("shutting down...")

	cancel()

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCtxCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}

	if err := database.Close(shutdownCtx); err != nil {
		logger.Error("failed to close database", "error", err)
	}

	signal.Stop(stop)
	close(stop)

	logger.Debug("done")
}
