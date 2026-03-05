//go:build !prod

package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	RegisterRouter()
}
