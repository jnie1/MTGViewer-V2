package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchUpdateLogs(c *gin.Context) {
	logs, err := transactions.FetchUpdateLogs()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, logs)
}

func fetchReportCards(c *gin.Context) {
	group := c.Param("group")
	group1, err := uuid.Parse(group)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	group2 := uuid.Nil
	if e, ok := c.GetQuery("e"); ok {
		group2, err = uuid.Parse(e)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	logs, err := fetchTransactionLogs(group1, group2)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	scryfallIds := transactions.GetScryfallIds(logs)
	cards, err := cards.FetchCollection(scryfallIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	reportCards, err := transactions.JoinReportCards(cards, logs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, reportCards)
}

func fetchTransactionLogs(group1 uuid.UUID, group2 uuid.UUID) ([]transactions.TransactionLogs, error) {
	if group2 == uuid.Nil {
		return transactions.FetchLogs(group1)
	}
	logRange, err := transactions.FetchLogRange(group1, group2)
	if err != nil {
		return nil, err
	}
	return transactions.FetchLogsFromRange(logRange)
}

func AddTransactionRoutes(router *gin.Engine) {
	group := router.Group("/logs")
	group.GET("", fetchUpdateLogs)
	group.GET("/:group", fetchReportCards)
}
