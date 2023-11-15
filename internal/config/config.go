package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	App
	Database
}

type App struct {
	GoEnv  string
	RunApp string
}

type Database struct {
	DataBaseUrl string
}

func getEnv(key, defaultKey string) string {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("[ERROR] loading .env file", err)
		}
	}

	var value = os.Getenv(key)

	if value == "" {
		return defaultKey
	}

	return value
}

func DefaultConfig() *Config {
	var c Config

	//KG: settings App
	c.GoEnv = getEnv("GO_ENV", "development")
	c.RunApp = getEnv("RUN_APP", "127.0.0.1:4000")

	// KG: settings Database
	c.DataBaseUrl = getEnv("DATABASE_URL", "")

	return &c
}
