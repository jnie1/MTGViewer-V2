package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	ClientOrigins []string
	Port          string
}

func Load() Config {
	originsEnv := os.Getenv("CLIENT_ORIGINS")
	if originsEnv == "" {
		log.Fatal("CLIENT_ORIGINS environment variable is not set")
	}

	var allowedOrigins []string
	for _, origin := range strings.Split(originsEnv, ",") {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			allowedOrigins = append(allowedOrigins, trimmedOrigin)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		ClientOrigins: allowedOrigins,
		Port:          port,
	}
}
