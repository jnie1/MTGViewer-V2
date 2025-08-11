package containers

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type Container struct {
	Name            string
	Capacity        int
	MarkForDeletion bool
}

type CardDeposit struct {
	ContainerId int
	ScryfallId  uuid.UUID
	Amount      int
}

type CardRequest struct {
	ScryfallId uuid.UUID
	Delta      int
}

type ContainerChanges struct {
	ContainerId int
	Requests    []CardRequest
}

type ContainerAllocation struct {
	ContainerId int
	Used        int
	MaxCapcity  int
}

func GetCardAmounts(deposits []CardDeposit, fullCards []cards.Card) []cards.CardAmount {
	amountMap := map[uuid.UUID]int{}

	for _, deposit := range deposits {
		amountMap[deposit.ScryfallId] = deposit.Amount
	}

	amounts := make([]cards.CardAmount, len(fullCards))

	for i, card := range fullCards {
		amount := amountMap[card.ScryfallId]
		amounts[i] = cards.CardAmount{Card: card, Amount: amount}
	}

	slices.SortFunc(amounts, func(a, b cards.CardAmount) int {
		return cmp.Compare(a.Amount, b.Amount)
	})

	return amounts
}

func GetScryfallIds(deposits []CardDeposit) []cards.ScryfallIdentifier {
	uniqIds := map[uuid.UUID]any{}

	for _, deposit := range deposits {
		uniqIds[deposit.ScryfallId] = nil
	}

	allIds := make([]cards.ScryfallIdentifier, len(uniqIds))
	i := 0

	for id := range uniqIds {
		allIds[i] = cards.ScryfallIdentifier{Id: id}
		i += 1
	}

	return allIds
}
