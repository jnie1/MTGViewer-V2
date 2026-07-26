package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
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

	deposits, err := containers.SearchDeposits(cardIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result, err := containers.JoinCardDeposits(cardPage.Cards, deposits)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func AddContainerRoutes(router *gin.Engine) {
	group := router.Group("/containers")
	group.GET("", fetchContainerPreviews)
	group.GET("/:container", fetchContainer)
	group.GET("/:container/cards", fetchContainerCards)
	group.GET("/cards", searchCards)
}
