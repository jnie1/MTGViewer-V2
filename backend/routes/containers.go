package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchContainerPreviews(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := containers.GetContainers(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchContainer(c *gin.Context) {
	var params struct {
		ContainerId int `uri:"container" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	result, err := containers.GetContainer(ctx, params.ContainerId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchContainerCards(c *gin.Context) {
	var params struct {
		ContainerId int `uri:"container" binding:"required"`
	}

	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	amounts, err := containers.GetAmounts(ctx, params.ContainerId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(amounts) == 0 {
		c.JSON(http.StatusOK, []cards.CardAmount{})
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	matches, err := cards.FetchCollection(ctx, scryfallIds...)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result, err := cards.JoinCardAmounts(matches, amounts)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func checkPrune(c *gin.Context) {
	var query struct {
		MaxCopies int     `form:"size" binding:"gte=0,required"`
		MinPrice  float64 `form:"price" binding:"gt=0.0,required"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	excess, err := containers.FindExcessDeposits(ctx, query.MaxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := containers.ToScryfallIds(excess)
	matches, err := cards.FetchCollection(ctx, scryfallIds...)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	prices, err := cards.FetchPrices(ctx, query.MinPrice, scryfallIds...)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes := containers.TranslatePrune(excess, matches, prices, query.MaxCopies, query.MinPrice)
	preview := containers.PreviewPrune(changes, matches, prices)

	c.JSON(http.StatusOK, preview)
}

func applyPrune(c *gin.Context) {
	var query struct {
		MaxCopies int     `form:"size" binding:"gte=0,required"`
		MinPrice  float64 `form:"price" binding:"gt=0.0,required"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	excess, err := containers.FindExcessDeposits(ctx, query.MaxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := containers.ToScryfallIds(excess)
	matches, err := cards.FetchCollection(ctx, scryfallIds...)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	prices, err := cards.FetchPrices(ctx, query.MinPrice, scryfallIds...)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes := containers.TranslatePrune(excess, matches, prices, query.MaxCopies, query.MinPrice)

	if err := containers.UpdateDeposits(ctx, changes); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := transactions.LogCollectionChanges(ctx, changes); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func AddContainerRoutes(router gin.IRouter) {
	group := router.Group("/containers")
	group.GET("", fetchContainerPreviews)
	group.GET("/:container", fetchContainer)
	group.GET("/:container/cards", fetchContainerCards)
	group.GET("/prune", checkPrune)
	group.POST("/prune", applyPrune)
}
