package routes

import (
	"mime/multipart"
	"net/http"
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

	ctx := c.Request.Context()
	result, err := cards.FetchCollection(ctx, scryfallIds)
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

	ctx := c.Request.Context()
	cardFound, err := cards.FetchCard(ctx, scryfallId)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	deposits, err := containers.SearchDepositsByOracleId(ctx, []uuid.UUID{cardFound.OracleId})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	result := containers.CardContainerMatch{
		Card:       cardFound,
		Containers: containers.MergeDespositsByContainer(deposits),
	}

	c.JSON(http.StatusOK, result)
}

func fetchRandomCard(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := cards.FetchRandomCard(ctx)

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

	ctx := c.Request.Context()
	requests, err := containers.ParseCardRequests(ctx, file)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	allocations, err := containers.GetAllocations(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	changes, err := containers.GetContainerAdditions(requests, allocations)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

func withdrawCards(c *gin.Context) {
	var withdrawals containers.ContainerWithdrawals
	if err := c.ShouldBind(&withdrawals); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	if err := containers.ResolveIdentifiers(ctx, withdrawals); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	scryfallIds := containers.FindScryfallIds(withdrawals)
	deposits, err := containers.SearchDeposits(ctx, scryfallIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	changes, err := containers.ValidateCardWithdrawals(withdrawals, deposits)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

	ctx := c.Request.Context()
	cardPage, err := cards.SearchCards(cardQuery, pageNum)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(cardPage.Cards) == 0 {
		c.JSON(http.StatusOK, cards.SearchCardPage{})
		return
	}

	oracleIds := make(uuid.UUIDs, len(cardPage.Cards))
	for i, card := range cardPage.Cards {
		oracleIds[i] = card.OracleId
	}

	cardMatches, err := containers.SearchDepositsByOracleId(ctx, oracleIds)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := cards.SearchCardPage{
		Page:    cardPage.Page,
		HasMore: cardPage.HasMore,
		Cards:   containers.FilterCards(cardPage.Cards, cardMatches),
	}

	c.JSON(http.StatusOK, result)
}

func refreshOracle(c *gin.Context) {
	ctx := c.Request.Context()

	ids, err := containers.FindMissingOracleIds(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(ids) == 0 {
		c.Status(http.StatusOK)
		return
	}

	matches, err := cards.FetchCollection(ctx, ids)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	objs := make([]cards.ScryfallOracleObj, len(matches))
	for i, card := range objs {
		objs[i] = cards.ScryfallOracleObj{ScryfallId: card.ScryfallId, OracleId: card.OracleId}
	}

	err = containers.UpdateOracleIds(ctx, objs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func AddCardRoutes(router gin.IRouter) {
	group := router.Group("/cards")
	group.GET("/", fetchCollection)
	group.GET("/:card", fetchCard)
	group.GET("/search", searchCards)
	group.GET("/random", fetchRandomCard)
	group.POST("/import", importCards)
	group.POST("/withdraw", withdrawCards)
	group.POST("/oracle", refreshOracle)
}
