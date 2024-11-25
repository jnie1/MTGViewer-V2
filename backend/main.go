package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/joho/godotenv"
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
	var jsonData map[string]any

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

func booksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, databaseStuffGoesHere)
}

func main() {
	// Loading environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.Open()
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer database.Close()

	// Initialize a Gin router
	r := gin.Default()

	// Define routes and associate them with handlers
	r.GET("/greet", greetHandler)
	r.GET("/user", userHandler)
	r.GET("/books", booksHandler)
	r.POST("/post", postHandler)

	// Start the server on port 8080
	r.Run(":8080")

}
