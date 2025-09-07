package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
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

	deposits, err := containers.GetDeposits(containerId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	scryfallIds := containers.GetScryfallIds(deposits)
	cards, err := cards.FetchCollection(scryfallIds)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cardAmounts, err := containers.GetCardAmounts(cards, deposits)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, cardAmounts)
}

func AddContainerRoutes(router *gin.Engine) {
	group := router.Group("/containers")
	group.GET("", fetchContainerPreviews)
	group.GET("/:container", fetchContainer)
	group.GET("/:container/cards", fetchContainerCards)
}
