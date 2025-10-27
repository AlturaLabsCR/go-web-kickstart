// Package config implements initialization logic for required app parameters
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Default values are initialized here, these will be used unless overwritten
	// by the Init() method

	AppTitle string = "My App Title"

	RootPrefix string = "/"

	// Depends on the RootPrefix, so, must be initialized after checking for any
	// overwrites of RootPrefix
	Assets string

	Production bool   = false
	Port       string = "8080"
	LogLevel   int    = 0 // -4:Debug 0:Info 4:Warn 8:Error
	dbDriver   string = "sqlite"
	dbConn     string = "./db.db"

	ServerSMTPUser string = "john@doe.com"
	ServerSMTPHost string
	ServerSMTPPort string
	ServerSMTPPass string
)

const (
	// You should use a prefix for any overwrites via env to avoid conflicts with
	// other programs
	envPrefix = "APP_"
	envPort   = envPrefix + "PORT"
	envProd   = envPrefix + "PROD"
	envLog    = envPrefix + "LOG_LEVEL"
	envDvr    = envPrefix + "DB_DRIVER"
	envCnn    = envPrefix + "DB_CONN"
	envRoot   = envPrefix + "ROOT_PREFIX"

	envSMTPHost = envPrefix + "SMTP_HOST"
	envSMTPPort = envPrefix + "SMTP_PORT"
	envSMTPPass = envPrefix + "SMTP_PASS"
)

func Init() {
	godotenv.Load()

	r := os.Getenv(envRoot)
	if r != "" {
		RootPrefix = r
	}

	Assets = RootPrefix + "assets/"

	Production = os.Getenv(envProd) == "1"

	p := os.Getenv(envPort)
	if p != "" {
		Port = p
	}

	ServerSMTPHost = os.Getenv(envSMTPHost)
	ServerSMTPPort = os.Getenv(envSMTPPort)
	ServerSMTPPass = os.Getenv(envSMTPPass)

	if ServerSMTPHost == "" || ServerSMTPPort == "" || ServerSMTPPass == "" {
		// Handle empty credentials
	}

	logLevelStr := os.Getenv(envLog)
	if logLevelStr != "" {
		var err error
		l, err := strconv.Atoi(logLevelStr)
		if err == nil {
			LogLevel = l
		}
	}

	dvr := os.Getenv(envDvr)
	if dvr != "" {
		dbDriver = dvr
	}

	conn := os.Getenv(envCnn)
	if conn != "" {
		dbConn = conn
	}
}
