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

func userHandler(c *gin.Context) {
	// Extracting the name parameter from the URL
	name := c.DefaultQuery("name", "Guest")

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, " + name,
	})
}

func postHandler(c *gin.Context) {
	var jsonData map[string]interface{}

	// Binding the incoming JSON request body to the map
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Received",
		"payload": jsonData,
	})
}

func main() {
	// Initialize a Gin router
	r := gin.Default()

	// Define routes and associate them with handlers
	r.GET("/greet", greetHandler)
	r.GET("/user", userHandler)
	r.POST("/post", postHandler)

	// Start the server on port 8080
	r.Run(":8080")
}
