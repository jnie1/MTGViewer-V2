//go:build prod

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handlers
func greetHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, welcome to the Go server!",
	})
}

func main() {
	RegisterRouter()
}
