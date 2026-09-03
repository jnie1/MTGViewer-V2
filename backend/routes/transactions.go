package routes

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/auth"
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

func fetchCardTransaction(c *gin.Context) {
	var params struct {
		Group uuid.UUID `uri:"group,parser=encoding.TextUnmarshaler" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx := c.Request.Context()
	result, err := transactions.GetTransaction(ctx, params.Group)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchCardLogs(c *gin.Context) {
	var params struct {
		Group1 uuid.UUID `uri:"group,parser=encoding.TextUnmarshaler" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var query struct {
		Group2 uuid.UUID `form:"e,parser=encoding.TextUnmarshaler"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	logs, err := getLogs(ctx, params.Group1, query.Group2)
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

func updateDescription(c *gin.Context) {
	var params struct {
		Group uuid.UUID `uri:"group,parser=encoding.TextUnmarshaler" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var req struct {
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	err := transactions.UpdateDescription(ctx, params.Group, req.Description)

	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func AddTransactionRoutes(router gin.IRouter) {
	group := router.Group("/logs")
	group.GET("", fetchCardTransactions)
	group.GET("/:group", fetchCardTransaction)
	group.GET("/:group/cards", fetchCardLogs)
	group.PUT("/:group/description", auth.IsAuthorized, auth.IsAdmin, updateDescription)
}
