//go:build prod

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/autotls"
)

func main() {
	r := CreateRouter()

	// this forces the port to be on port 80
	if err := autotls.Run(r, os.Getenv("API_DOMAIN")); err != nil {
		log.Fatal(err)
	}
}
