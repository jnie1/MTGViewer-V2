//go:build !prod

package main

import (
	"log"

	"github.com/jnie1/MTGViewer-V2/config"
	"github.com/joho/godotenv"
)

func main() {
	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()
	RegisterRouter(cfg)
}
