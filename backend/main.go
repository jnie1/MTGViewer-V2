package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/auth"
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

type bookInfo struct {
	Name       string `json:"name"`
	AuthorName string `json:"authorName"`
}

func getBooks(c *gin.Context) {
	books := []bookInfo{}

	db := database.Instance()
	row, err := db.Query(`
		SELECT books.name, author.name
		FROM books
		JOIN author ON books.author_id = author.id`)

	if err != nil {
		c.JSON(http.StatusOK, books)
		return
	}

	for row.Next() {
		book := bookInfo{}
		if err = row.Scan(&book.Name, &book.AuthorName); err == nil {
			books = append(books, book)
		}
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	book := bookInfo{}

	db := database.Instance()

	row := db.QueryRow(`
		SELECT books.name, author.name
		FROM books
		JOIN author ON books.author_id = author.id
		WHERE id = $1`, bookId)

	err = row.Scan(&book.Name, &book.AuthorName)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
}

func addBook(c *gin.Context) {
	var bookInfo bookInfo

	if err := c.ShouldBind(&bookInfo); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	db := database.Instance()
	authorId, err := getAuthorId(db, bookInfo.AuthorName)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = db.Exec("INSERT INTO books (name, author_id) VALUES ($1, $2)", bookInfo.Name, authorId)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, bookInfo)
}

func getAuthorId(db *sql.DB, authorName string) (int, error) {
	var authorId int

	row := db.QueryRow("SELECT id FROM author WHERE name = $1", authorName)
	err := row.Scan(&authorId)

	if err == nil {
		return authorId, nil
	}

	row = db.QueryRow("INSERT INTO author (name) VALUES ($1) RETURNING id", authorName)
	err = row.Scan(&authorId)

	return authorId, err
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
	r.GET("/books/:id", getBook)
	r.GET("/books", getBooks)
	r.POST("/books", addBook)
	r.POST("/post", postHandler)

	authorized := r.Group("", auth.IsAuthorized)

	authorized.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"secret": "some secret"})
	})

	// Start the server on port 8080
	r.Run(":8080")

}
