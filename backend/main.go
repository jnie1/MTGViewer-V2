package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

func booksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, databaseStuffGoesHere)
}

func main() {
	// Initialize a Gin router
	r := gin.Default()

	// Loading environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieving environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Database connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	// Opening connection to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	// Pinging databse connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	// Define routes and associate them with handlers
	r.GET("/greet", greetHandler)
	r.GET("/user", userHandler)
	r.GET("/books", booksHandler)
	r.POST("/post", postHandler)

	// Start the server on port 8080
	r.Run(":8080")

}
