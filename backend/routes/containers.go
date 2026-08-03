package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchContainerPreviews(c *gin.Context) {
	result, err := containers.GetContainers()
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

	result, err := containers.GetContainer(containerId)
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

	amounts, err := containers.GetAmounts(containerId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(amounts) == 0 {
		c.JSON(http.StatusOK, []cards.CardAmount{})
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	matches, err := cards.FetchCollection(scryfallIds)
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

	amounts, err := containers.FindExcessAmounts(maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	matches, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := cards.FindCheapCards(matches, amounts, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
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

	amounts, err := containers.FindExcessAmounts(maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.ToScryfallIds(amounts)
	results, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cheapCards, err := cards.FindCheapCards(results, amounts, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cheapIds := make(uuid.UUIDs, len(cheapCards))
	for i, card := range cheapCards {
		cheapIds[i] = card.ScryfallId
	}

	matches, err := containers.SearchDeposits(cheapIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes := containers.TranslatePrune(matches, maxCopies)

	if err := containers.UpdateDeposits(changes); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := transactions.LogCollectionChanges(changes); err != nil {
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
	group.GET("/cards", searchCards)
	group.GET("/prune", checkPrune)
	group.POST("/prune", applyPrune)
}
