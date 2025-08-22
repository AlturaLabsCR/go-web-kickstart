package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/config"
	"app/handlers"
	"app/router"
)

const (
	defaultPort = "8080"
)

func main() {
	var (
		port       = os.Getenv("APP_PORT")
		production = os.Getenv("APP_PROD") == "1"
		logLevel   = os.Getenv("APP_LOG_LEVEL")
		dbDriver   = os.Getenv("APP_DB_DRIVER")
		dbConn     = os.Getenv("APP_DB_CONN")
	)

	if port == "" {
		port = defaultPort
	}

	logger, err := config.InitLogger(production, logLevel)
	if err != nil {
		print("failed logger initialization: %v\n", err)
		os.Exit(1)
	}

	database, err := config.InitDB(dbDriver, dbConn)
	if err != nil {
		print("failed database initialization: %v\n", err)
		os.Exit(1)
	}

	handler := handlers.New(handlers.HandlerParams{
		Production: production,
		DB:         database,
		Logger:     logger,
	})

	routes := router.Routes(handler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", "address", ":"+port)

		if err := http.ListenAndServe(":"+port, routes); err != nil {
			logger.Error("failed to start server", "port", port, "error", err)
			os.Exit(1)
		}
	}()

	<-stop

	logger.Info("shutting down...")
}
