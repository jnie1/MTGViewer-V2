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

func main() {
	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.Open(); err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer database.Close()

	// Initialize a Gin router
	r := gin.Default()
	r.Use(auth.AddCors)

	// Define routes and associate them with handlers
	routes.AddUserRoutes(r)
	routes.AddCardRoutes(r)
	routes.AddContainerRoutes(r)
	routes.AddTransactionRoutes(r)

	authorized := r.Group("", auth.IsAuthorized)

	authorized.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"secret": "some secret"})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
