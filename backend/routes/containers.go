package routes

import (
	"net/http"
	"strconv"

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
	id := c.Param("container")
	containerId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	result, err := containers.GetContainer(ctx, containerId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchContainerCards(c *gin.Context) {
	id := c.Param("container")
	containerId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	amounts, err := containers.GetAmounts(ctx, containerId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(amounts) == 0 {
		c.JSON(http.StatusOK, []cards.CardAmount{})
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	matches, err := cards.FetchCollection(ctx, scryfallIds)

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
	size := c.Query("size")
	maxCopies, err := strconv.Atoi(size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if maxCopies <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	price := c.Query("price")
	maxPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if maxPrice <= 0.0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	amounts, err := containers.FindExcessAmounts(ctx, maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	matches, err := cards.FetchCollection(ctx, scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	prices, err := cards.FetchPrices(ctx, scryfallIds, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cheapCards := containers.FindCheapCandidates(matches, prices, maxPrice)
	deposits, err := containers.SearchDeposits(ctx, cheapCards)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes := containers.TranslatePrune(deposits, maxCopies)
	preview, err := containers.PreviewPrune(changes, matches, prices)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, preview)
}

func applyPrune(c *gin.Context) {
	size := c.Query("size")
	maxCopies, err := strconv.Atoi(size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if maxCopies <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	price := c.Query("price")
	maxPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if maxPrice <= 0.0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	amounts, err := containers.FindExcessAmounts(ctx, maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)

	matches, err := cards.FetchCollection(ctx, scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	prices, err := cards.FetchPrices(ctx, scryfallIds, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cheapCards := containers.FindCheapCandidates(matches, prices, maxPrice)
	deposits, err := containers.SearchDeposits(ctx, cheapCards)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes := containers.TranslatePrune(deposits, maxCopies)

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
