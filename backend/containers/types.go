package containers

import (
	"cmp"
	"slices"
	"strings"

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
	MaxCapacity int
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

type csvHeaderPositions struct {
	Name            int
	ScryfallId      int
	MultiverseId    int
	SetCode         int
	CollectorNumber int
	Quantity        int
}

func (positions *csvHeaderPositions) hasValidPosition() bool {
	if positions == nil {
		return false
	}

	if positions.Quantity == -1 {
		return false
	}

	if positions.ScryfallId > -1 {
		return true
	}

	if positions.MultiverseId > -1 {
		return true
	}

	if positions.SetCode > -1 && positions.CollectorNumber > -1 {
		return true
	}

	if positions.SetCode > -1 && positions.Name > -1 {
		return true
	}

	return false
}

func getHeaderPositions(header []string) csvHeaderPositions {
	return csvHeaderPositions{
		Name:            getHeaderIndex(header, "name"),
		ScryfallId:      getHeaderIndex(header, "scryfall id", "scryfallid"),
		MultiverseId:    getHeaderIndex(header, "multiverse id", "multiverseid"),
		SetCode:         getHeaderIndex(header, "set code", "setcode"),
		CollectorNumber: getHeaderIndex(header, "collector number", "collectornumber"),
		Quantity:        getHeaderIndex(header, "quantity"),
	}
}

func getHeaderIndex(header []string, targetNames ...string) int {
	searchTargetNames := func(column string) bool {
		matchIndex := slices.IndexFunc(targetNames, func(target string) bool {
			return strings.EqualFold(column, target)
		})
		return matchIndex > -1
	}
	return slices.IndexFunc(header, searchTargetNames)
}
