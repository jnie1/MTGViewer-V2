package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchCardTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := transactions.GetTransactions(ctx)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchCardLogs(c *gin.Context) {
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

	ctx := c.Request.Context()
	logs, err := getLogs(ctx, group1, group2)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logs = transactions.MergeLogs(logs)
	scryfallIds := transactions.ToScryfallIds(logs)

	matches, err := cards.FetchCollection(ctx, scryfallIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result, err := transactions.JoinCardLogs(matches, logs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func getLogs(ctx context.Context, group1, group2 uuid.UUID) ([]transactions.CardLogPreview, error) {
	if group2 == uuid.Nil {
		return transactions.GetLogs(ctx, group1)
	}
	logRange, err := transactions.GetTimeRange(ctx, group1, group2)
	if err != nil {
		return nil, err
	}
	return transactions.GetLogsFromRange(ctx, logRange)
}

func AddTransactionRoutes(router gin.IRouter) {
	group := router.Group("/logs")
	group.GET("", fetchCardTransactions)
	group.GET("/:group", fetchCardLogs)
}
