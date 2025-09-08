package containers

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type ContainerPreview struct {
	ContainerId int    `json:"containerId"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
}

type Container struct {
	Name      string `json:"name"`
	Used      int    `json:"used"`
	Capacity  int    `json:"capacity"`
	IsDeleted bool   `json:"isDeleted"`
}

type CardDepositAmount struct {
	ContainerId   int    `json:"containerId"`
	ContainerName string `json:"containerName"`
	cards.CardAmount
}

type CardDepositPreview struct {
	ScryfallId uuid.UUID
	Amount     int
}

type CardDeposit struct {
	ContainerId   int
	ContainerName string
	CardDepositPreview
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

func (allocation ContainerAllocation) Remaining() int {
	return allocation.MaxCapacity - allocation.Used
}

func GetCardAmounts(fullCards []cards.Card, deposits []CardDepositPreview) ([]cards.CardAmount, error) {
	amounts := make([]cards.CardAmount, len(deposits))
	cardMap := make(map[uuid.UUID]cards.Card, len(fullCards))

	for _, card := range fullCards {
		cardMap[card.ScryfallId] = card
	}

	for i, deposit := range deposits {
		card, ok := cardMap[deposit.ScryfallId]
		if !ok {
			return nil, fmt.Errorf("cannot resolve card id %s", deposit.ScryfallId)
		}
		amounts[i] = cards.CardAmount{
			Card:   card,
			Amount: deposit.Amount,
		}
	}

	slices.SortFunc(amounts, func(a, b cards.CardAmount) int {
		return cmp.Compare(a.Amount, b.Amount)
	})

	return amounts, nil
}

func GetCardDepositAmounts(fullCards []cards.Card, deposits []CardDeposit) ([]CardDepositAmount, error) {
	depositAmounts := make([]CardDepositAmount, len(deposits))
	cardMap := make(map[uuid.UUID]cards.Card, len(fullCards))

	for _, card := range fullCards {
		cardMap[card.ScryfallId] = card
	}

	for i, deposit := range deposits {
		card, ok := cardMap[deposit.ScryfallId]
		if !ok {
			return nil, fmt.Errorf("cannot resolve card id %s", deposit.ScryfallId)
		}
		amount := cards.CardAmount{Card: card, Amount: deposit.Amount}
		depositAmounts[i] = CardDepositAmount{CardAmount: amount, ContainerId: deposit.ContainerId, ContainerName: deposit.ContainerName}
	}

	slices.SortFunc(depositAmounts, func(a, b CardDepositAmount) int {
		return cmp.Compare(a.Amount, b.Amount)
	})

	return depositAmounts, nil
}

func GetScryfallIds(deposits []CardDepositPreview) []cards.ScryfallIdentifier {
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

func (positions *csvHeaderPositions) Valid() bool {
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
