package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchTransactionLogs(c *gin.Context) {
	listOfLogs, err := transactions.FetchLogs()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	scryfallIds := transactions.GetScryfallIds(listOfLogs)
	cards, err := cards.FetchCollection(scryfallIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	listOfReportCards, err := transactions.JoinReportCards(cards, listOfLogs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, listOfReportCards)
}

func AddTransactionRoutes(router *gin.Engine) {
	group := router.Group("/transactions")
	group.GET("/logs", fetchTransactionLogs)
}
