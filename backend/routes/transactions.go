package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
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
	group1, err := uuid.Parse(c.Param("group"))
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

	containerIds := transactions.ToContainerIds(logs)
	boxes, err := containers.GetContainerPreviews(ctx, containerIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	scryfallIds := transactions.ToScryfallIds(logs)
	cardMatches, err := cards.FetchCollection(ctx, scryfallIds...)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	transfers, err := transactions.MergeLogs(logs, boxes, cardMatches)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
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
