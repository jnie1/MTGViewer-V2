package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/transactions"
	"net/http"
)

func fetchTransactionLogs(c *gin.Context) {
	listOfLogs, err := transactions.FetchLogs()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, listOfLogs)
}

func AddTransactionRoutes(router *gin.Engine) {
	group := router.Group("/transactions")
	group.GET("/logs", fetchTransactionLogs)
}
