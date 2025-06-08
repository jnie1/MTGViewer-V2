package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
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
	scryfallId, err := uuid.Parse((cardId))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	card, err := cards.FetchCard(scryfallId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, card)
}

func fetchCollection(c *gin.Context) {
	cardIds := c.QueryArray("cards")
	scryfallIds := make([]uuid.UUID, len(cardIds))

	for i, cardId := range cardIds {
		id, err := uuid.Parse(cardId)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		scryfallIds[i] = id
	}

	cards, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func AddCardRoutes(router *gin.Engine) {
	group := router.Group("/cards")
	group.GET("/scryfall", fetchRandomCard)
	group.GET("/:card", fetchCard)
	group.GET("/collection", fetchCollection)
}
