package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	transactions "github.com/jnie1/MTGViewer-V2/transaction"
)

func fetchTransactionLogs(c *gin.Context) {
	listOfLogs, err := transactions.FetchLogs()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, listOfLogs)
}

func fetchLogsRoutes(router *gin.Engine) {
	group := router.Group("/transactions")
	group.GET("/logs", fetchTransactionLogs)
}
