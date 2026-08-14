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

type ContainerEntry struct {
	Name      string `json:"name"`
	Used      int    `json:"used"`
	Capacity  int    `json:"capacity"`
	IsDeleted bool   `json:"isDeleted"`
}

type Container struct {
	ContainerAllocation
	Name string `json:"name"`
}

type ContainerPreview struct {
	ContainerId int    `json:"containerId"`
	Name        string `json:"name"`
}

func (container *ContainerPreview) Container() ContainerPreview {
	if container == nil {
		return ContainerPreview{}
	}
	return *container
}

type ContainerDeposit struct {
	ContainerPreview
	Amount int `json:"amount"`
}

type ContainerDelta struct {
	ContainerId int
	Delta       int
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

type CardContainerMatch struct {
	Card       cards.Card         `json:"card"`
	Containers []ContainerDeposit `json:"containers"`
}

type CardPrunePreview struct {
	Total int                     `json:"total"`
	Cards []cards.CardPriceAmount `json:"cards"`
}

type CardRequest struct {
	ScryfallId uuid.UUID
	OracleId   uuid.UUID
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
	Card   any `json:"card"`
	Amount int `json:"amount"`
}

func (id *CardIdentifierAmount) UnmarshalJSON(data []byte) error {
	var obj struct {
		Card   map[string]any `json:"card"`
		Amount int            `json:"amount"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	card, err := cards.FromObj(obj.Card)
	if err != nil {
		return err
	}

	id.Card = card
	id.Amount = obj.Amount
	return nil
}

type ContainerAllocation struct {
	ContainerId int `json:"containerId"`
	Used        int `json:"used"`
	Capacity    int `json:"capacity"`
}

func (allocation ContainerAllocation) Remaining() int {
	return allocation.Capacity - allocation.Used
}

func CompareRemaining(a, b ContainerAllocation) int {
	return cmp.Compare(a.Remaining(), b.Remaining())
}

func mergeCardRequestsChecked(requests []CardRequest) ([]CardRequest, bool) {
	cardCounter := map[uuid.UUID]int{}
	oracleIds := map[uuid.UUID]uuid.UUID{}
	for _, request := range requests {
		cardCounter[request.ScryfallId] = cardCounter[request.ScryfallId] + request.Delta
		oracleIds[request.ScryfallId] = request.OracleId
	}

	if len(cardCounter) == len(requests) {
		// no requests were merged, so just use original array
		return requests, false
	}

	mergedRequests := make([]CardRequest, len(cardCounter))
	i := 0

	for cardId, delta := range cardCounter {
		mergedRequests[i] = CardRequest{cardId, oracleIds[cardId], delta}
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

func MergeDespositsByContainer(cardDeposits []CardDepositPreview) []ContainerDeposit {
	boxAmounts := map[int]int{}
	for _, deposit := range cardDeposits {
		boxAmounts[deposit.ContainerId] = boxAmounts[deposit.ContainerId] + deposit.Amount
	}

	boxDeposits := make([]ContainerDeposit, len(boxAmounts))

	for _, card := range cardDeposits {
		if amount, ok := boxAmounts[card.ContainerId]; ok {
			idx := len(boxDeposits) - len(boxAmounts)
			deposit := ContainerDeposit{
				ContainerPreview{card.ContainerId, card.ContainerName},
				amount,
			}
			boxDeposits[idx] = deposit
			delete(boxAmounts, card.ContainerId)
		}
	}

	return boxDeposits
}

func FilterCards(fullCards []cards.Card, deposits []CardDepositPreview) []cards.Card {
	oracleIds := map[uuid.UUID]any{}
	for _, deposit := range deposits {
		oracleIds[deposit.OracleId] = nil
	}

	matches := make([]cards.Card, len(oracleIds))
	i := 0

	for _, card := range fullCards {
		if _, ok := oracleIds[card.OracleId]; ok {
			matches[i] = card
			i += 1
		}
	}

	slices.SortFunc(matches, func(a, b cards.Card) int {
		nameCompare := strings.Compare(a.Name, b.Name)
		if nameCompare == 0 {
			return strings.Compare(a.Set, b.Set)
		}
		return nameCompare
	})

	return matches
}

func ToScryfallIds(deposits []CardDepositPreview) uuid.UUIDs {
	uniqIds := map[uuid.UUID]any{}

	for _, deposit := range deposits {
		uniqIds[deposit.ScryfallId] = nil
	}

	ids := make(uuid.UUIDs, len(uniqIds))
	i := 0

	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
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
