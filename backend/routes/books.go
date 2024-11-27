package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/books"
)

func getBooks(c *gin.Context) {
	books, err := books.GetBooks()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
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

	book, err := books.GetBook(bookId)

	if err == sql.ErrNoRows {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func addBook(c *gin.Context) {
	var book books.BookData

	err := c.ShouldBind(&book)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = books.AddBook(book)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, book)
}

func AddBookRoutes(router *gin.Engine) {
	group := router.Group("/books")
	group.GET("/:id", getBook)
	group.GET("/", getBooks)
	group.POST("/", addBook)
}
