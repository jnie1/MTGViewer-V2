package routes

import (
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
	"github.com/jnie1/MTGViewer-V2/transactions"
)

func fetchCollection(c *gin.Context) {
	ids := c.QueryArray("cards")

	if len(ids) == 0 {
		c.JSON(http.StatusOK, []cards.Card{})
		return
	}

	scryfallIds, err := cards.ParseScryfallIds(ids)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := cards.FetchCollection(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func fetchCard(c *gin.Context) {
	cardId := c.Param("card")

	scryfallId, err := uuid.Parse(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cardFound, err := cards.FetchCard(cards.ScryfallIdentifier{Id: scryfallId})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	allPrints, err := cards.SearchCards(cardFound.Name, 1)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	scryfallIds := make([]uuid.UUID, len(allPrints.Cards))
	for i, card := range allPrints.Cards {
		scryfallIds[i] = card.ScryfallId
	}

	// find scryfall id in containers
	containerResult, err := containers.SearchDeposits(scryfallIds)

	mergedBoxAmount := make(map[int]containers.CardDepositPreview, len(containerResult))

	for _, containerResult := range containerResult {
		if existing, ok := mergedBoxAmount[containerResult.ContainerId]; ok {
			existing.Amount = existing.Amount + containerResult.Amount
			mergedBoxAmount[containerResult.ContainerId] = existing
		} else {
			mergedBoxAmount[containerResult.ContainerId] = containerResult
		}
	}

	mergedResult := make([]containers.CardDepositPreview, 0, len(mergedBoxAmount))
	for _, containerResult := range mergedBoxAmount {
		mergedResult = append(mergedResult, containerResult)
	}

	sort.Slice(mergedResult, func(i, j int) bool {
		return mergedResult[i].ContainerId < mergedResult[j].ContainerId
	})

	sortedMergedResult := mergedResult

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	result := containers.CardContainerMatch{
		Card:       cardFound,
		Containers: sortedMergedResult,
	}
	c.JSON(http.StatusOK, result)
}

func fetchRandomCard(c *gin.Context) {
	result, err := cards.FetchRandomCard()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
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

	changes, err := containers.GetContainerAdditions(requests, allocations)
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

func withdrawCards(c *gin.Context) {
	var withdrawals containers.ContainerWithdrawals
	if err := c.ShouldBind(&withdrawals); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := containers.ResolveExtraIdentifiers(withdrawals); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := containers.FindScryfallIds(withdrawals)
	deposits, err := containers.SearchDeposits(scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	changes, err := containers.ValidateCardWithdrawals(withdrawals, deposits)
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

func searchCards(c *gin.Context) {
	cardQuery := c.Query("q")
	cardPages := c.Query("page")
	if cardPages == "" {
		cardPages = "1"
	}

	pageNum, err := strconv.Atoi(cardPages)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cardPage, err := cards.SearchCards(cardQuery, pageNum)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(cardPage.Cards) == 0 {
		c.JSON(http.StatusOK, cards.SearchCardPage{})
		return
	}

	cardIds := make(uuid.UUIDs, len(cardPage.Cards))
	for i, card := range cardPage.Cards {
		cardIds[i] = card.ScryfallId
	}

	cardMatches, err := containers.MatchCards(cardIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	filteredCards, err := cards.FilterCards(cardPage.Cards, cardMatches)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := cards.SearchCardPage{
		TotalCards: cardPage.TotalCards,
		Page:       cardPage.Page,
		HasMore:    cardPage.HasMore,
		Cards:      filteredCards,
	}

	c.JSON(http.StatusOK, result)
}

func AddCardRoutes(router gin.IRouter) {
	group := router.Group("/cards")
	group.GET("/", fetchCollection)
	group.GET("/:card", fetchCard)
	group.GET("/search", searchCards)
	group.GET("/random", fetchRandomCard)
	group.POST("/import", importCards)
	group.POST("/withdraw", withdrawCards)
	group.GET("/test", func(c *gin.Context) {
		result, err := cards.TestFiltering()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
