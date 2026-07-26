package containers

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type ContainerName struct {
	ContainerId int    `json:"containerId"`
	Name        string `json:"name"`
}

func (container *ContainerName) Container() ContainerName {
	if container == nil {
		return ContainerName{}
	}
	return *container
}

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

type CardDeposit struct {
	ContainerId   int    `json:"containerId"`
	ContainerName string `json:"containerName"`
	cards.CardAmount
}

type CardDepositPreview struct {
	ContainerId   int
	ContainerName string
	cards.CardAmountPreview
}

type CardRequest struct {
	ScryfallId uuid.UUID
	Delta      int
}

type ContainerCard struct {
	ContainerId int
	ScryfallId  uuid.UUID
}

type ContainerChanges struct {
	ContainerId int
	Requests    []CardRequest
}

type ContainerWithdrawals map[int][]CardIdentifierAmount

type CardIdentifierAmount struct {
	Card   CardIdentifier `json:"card"`
	Amount int            `json:"amount"`
}

func (id *CardIdentifierAmount) UnmarshalJSON(data []byte) error {
	var obj struct {
		Card   map[string]any `json:"card"`
		Amount int            `json:"amount"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	card, err := FromObj(obj.Card)
	if err != nil {
		return err
	}

	id.Card = card
	id.Amount = obj.Amount
	return nil
}

type ContainerAllocation struct {
	ContainerId int
	Used        int
	MaxCapacity int
}

func (allocation ContainerAllocation) Remaining() int {
	return allocation.MaxCapacity - allocation.Used
}

func CompareRemaining(a, b ContainerAllocation) int {
	return cmp.Compare(a.Remaining(), b.Remaining())
}

func mergeCardRequestsChecked(requests []CardRequest) ([]CardRequest, bool) {
	cardCounter := map[uuid.UUID]int{}
	for _, request := range requests {
		cardCounter[request.ScryfallId] = cardCounter[request.ScryfallId] + request.Delta
	}

	if len(cardCounter) == len(requests) {
		// no requests were merged, so just use original array
		return requests, false
	}

	mergedRequests := make([]CardRequest, len(cardCounter))
	i := 0

	for cardId, delta := range cardCounter {
		mergedRequests[i] = CardRequest{cardId, delta}
		i += 1
	}

	return mergedRequests, true
}

func MergeCardRequests(requests []CardRequest) []CardRequest {
	merged, _ := mergeCardRequestsChecked(requests)
	return merged
}

func MergeContainerChanges(changes []ContainerChanges) []ContainerChanges {
	requestsById := map[int][]CardRequest{}
	for _, change := range changes {
		requestsById[change.ContainerId] = append(requestsById[change.ContainerId], change.Requests...)
	}

	hasChanges := false
	mergedChanges := []ContainerChanges{}

	for containerId, requests := range requestsById {
		mergedRequests, requestsChanged := mergeCardRequestsChecked(requests)
		if requestsChanged {
			hasChanges = true
		}
		mergedChanges = append(mergedChanges, ContainerChanges{containerId, mergedRequests})
	}

	if !hasChanges && len(mergedChanges) == len(changes) {
		return changes
	}
	return mergedChanges
}

func JoinCardDeposits(fullCards []cards.Card, deposits []CardDepositPreview) ([]CardDeposit, error) {
	depositAmounts := make([]CardDeposit, len(deposits))
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
		depositAmounts[i] = CardDeposit{CardAmount: amount, ContainerId: deposit.ContainerId, ContainerName: deposit.ContainerName}
	}

	slices.SortFunc(depositAmounts, func(a, b CardDeposit) int {
		nameCompare := strings.Compare(a.Name, b.Name)
		if nameCompare == 0 {
			return strings.Compare(a.Set, b.Set)
		}
		return nameCompare
	})

	return depositAmounts, nil
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
