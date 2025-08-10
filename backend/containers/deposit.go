package containers

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func GetCardAmounts(deposits []CardDeposit, fullCards []cards.Card) []cards.CardAmount {
	amountMap := map[uuid.UUID]int{}
	for _, deposit := range deposits {
		amountMap[deposit.ScryfallId] = deposit.Amount
	}

	amounts := []cards.CardAmount{}
	for _, card := range fullCards {
		amount := amountMap[card.ScryfallId]
		amounts = append(amounts, cards.CardAmount{Card: card, Amount: amount})
	}

	slices.SortFunc(amounts, func(a, b cards.CardAmount) int {
		return cmp.Compare(a.Amount, b.Amount)
	})

	return amounts
}

func GetScryfallIds(deposits []CardDeposit) uuid.UUIDs {
	uniqIds := map[uuid.UUID]any{}
	for _, deposit := range deposits {
		uniqIds[deposit.ScryfallId] = nil
	}

	allIds := uuid.UUIDs{}
	for id := range uniqIds {
		allIds = append(allIds, id)
	}

	return allIds
}
