package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/books"
)

func getBooks(c *gin.Context) {
	books, _ := books.GetBooks()
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	bookId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	book, err := books.GetBook(bookId)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
}

func addBook(c *gin.Context) {
	var book books.BookData

	if err := c.ShouldBind(&book); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	err := books.AddBook(book)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, book)
}

func AddBookRoutes(router *gin.Engine) {
	router.GET("/books/:id", getBook)
	router.GET("/books", getBooks)
	router.POST("/books", addBook)
}
