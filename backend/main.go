package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
	"github.com/jnie1/MTGViewer-V2/database"
	"github.com/jnie1/MTGViewer-V2/routes"
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
	routes.AddBookRoutes(r)
	r.GET("/greet", greetHandler)
	r.GET("/user", userHandler)

	authorized := r.Group("", auth.IsAuthorized)

	authorized.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"secret": "some secret"})
	})

	// Start the server on port 8080
	r.Run(":8080")

}
