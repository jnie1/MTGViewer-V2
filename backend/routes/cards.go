package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func fetchCollection(c *gin.Context) {
	cardIds := c.QueryArray("cards")
	cards, err := cards.FetchCollection(cardIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func AddCardRoutes(router *gin.Engine) {
	group := router.Group("/cards")
	group.GET("/scryfall", fetchRandomCard)
	group.GET("/collection", fetchCollection)
}
