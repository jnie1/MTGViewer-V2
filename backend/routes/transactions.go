package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchUpdateLogs(c *gin.Context) {
	listOfLogs, err := transactions.FetchUpdateLogs()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, listOfLogs)
}

func fetchTransactionLogs(c *gin.Context) {
	group := c.Param("group")

	groupId, err := uuid.Parse(group)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	listOfLogs, err := transactions.FetchLogs(groupId)
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

func fetchTransactionRange(c *gin.Context) {
	g1 := c.Param("group1")
	g2 := c.Param("group2")

	group1, err := uuid.Parse(g1)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	group2, err := uuid.Parse(g2)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	logRange, err := transactions.FetchLogRange(group1, group2)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	listOfLogs, err := transactions.FetchLogsFromRange(logRange)
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
	group := router.Group("/logs")
	group.GET("", fetchUpdateLogs)
	group.GET("/:group", fetchTransactionLogs)
	group.GET("/:group1/:group2", fetchTransactionRange)
}
