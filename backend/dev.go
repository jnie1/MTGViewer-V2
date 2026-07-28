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

	r := CreateRouter()

	// // Start the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
