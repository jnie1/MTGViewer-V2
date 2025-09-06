package routes

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchRandomCard(c *gin.Context) {
	card, err := cards.FetchRandomCard()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, card)
}

func fetchCard(c *gin.Context) {
	cardId := c.Param("card")

	scryfallId, err := uuid.Parse(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card, err := cards.FetchCard(cards.ScryfallIdentifier{Id: scryfallId})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, card)
}

func fetchCollection(c *gin.Context) {
	ids := c.QueryArray("cards")

	scryfallIds, err := cards.ParseScryfallIds(ids)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cards, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func importCards(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if file.Size >= 5_000_000 {
		c.AbortWithError(http.StatusBadRequest, multipart.ErrMessageTooLarge)
		return
	}

	requests, err := containers.ParseCardRequests(file)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	allocations, err := containers.GetAllocations()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	changes, err := containers.GetContainerChanges(requests, allocations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
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

func AddCardRoutes(router *gin.Engine) {
	group := router.Group("/cards")
	group.GET("/scryfall", fetchRandomCard)
	group.GET("/:card", fetchCard)
	group.GET("/collection", fetchCollection)
	group.POST("/import", importCards)
}
