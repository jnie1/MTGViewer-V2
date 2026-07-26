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
	containers, err := containers.GetContainers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, containers)
}

func fetchContainer(c *gin.Context) {
	id := c.Param("container")
	containerId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	container, err := containers.GetContainer(containerId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, container)
}

func fetchContainerCards(c *gin.Context) {
	id := c.Param("container")
	containerId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	previews, err := containers.GetAmountPreviews(containerId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(previews) == 0 {
		c.JSON(http.StatusOK, []cards.CardAmount{})
		return
	}

	scryfallIds := cards.GetScryfallIds(previews)
	matches, err := cards.FetchCollection(scryfallIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	amounts, err := cards.GetCardAmounts(matches, previews)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, amounts)
}

func searchCards(c *gin.Context) {
	cardQuery := c.Query("q")

	cardPage, err := cards.SearchCards(cardQuery, 1)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(cardPage.Cards) == 0 {
		c.JSON(http.StatusOK, []containers.CardDeposit{})
		return
	}

	cardIds := make(uuid.UUIDs, len(cardPage.Cards))
	for i, card := range cardPage.Cards {
		cardIds[i] = card.ScryfallId
	}

	previews, err := containers.SearchCards(cardIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	amounts, err := containers.GetCardDeposits(cardPage.Cards, previews)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, amounts)
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

	previews, err := containers.FindCardsAboveCount(maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.GetScryfallIds(previews)
	results, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pruningCards, err := cards.KeepBelowPrice(results, previews, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pruningCards)
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

	previews, err := containers.FindCardsAboveCount(maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := cards.GetScryfallIds(previews)
	results, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pruningCards, err := cards.KeepBelowPrice(results, previews, maxPrice)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pruningIds := make(uuid.UUIDs, len(pruningCards))
	for i, card := range pruningCards {
		pruningIds[i] = card.ScryfallId
	}

	matches, err := containers.SearchCards(pruningIds)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	changes, err := containers.TranslatePrune(matches, maxCopies)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

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

func AddContainerRoutes(router *gin.Engine) {
	group := router.Group("/containers")
	group.GET("", fetchContainerPreviews)
	group.GET("/:container", fetchContainer)
	group.GET("/:container/cards", fetchContainerCards)
	group.GET("/cards", searchCards)
	group.GET("/prune", checkPrune)
	group.POST("/prune", applyPrune)
}
