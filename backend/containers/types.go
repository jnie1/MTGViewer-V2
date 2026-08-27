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
	SortOrder   int    `json:"sortOrder"`
}

type ContainerDeposit struct {
	ContainerId int                     `json:"containerId"`
	Name        string                  `json:"name"`
	Amount      int                     `json:"amount"`
	Prints      []cards.CardImageAmount `json:"prints"`
}

type ContainerDelta struct {
	ContainerId int
	Delta       int
}

type CardDeposit struct {
	ContainerId int `json:"containerId"`
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

	var mergedChanges []ContainerChanges
	hasChanges := false

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

func JoinContainerDeposits(boxes []ContainerPreview, deposits []CardDeposit) ([]ContainerDeposit, error) {
	boxPrints := map[int][]cards.CardImageAmount{}

	for _, deposit := range deposits {
		image, err := cards.ImageURLs(deposit.ScryfallId)
		if err != nil {
			return nil, err
		}

		cardImage := cards.CardImageAmount{
			ScryfallId: deposit.ScryfallId,
			Images:     image,
			Amount:     deposit.Amount,
		}

		boxPrints[deposit.ContainerId] = append(boxPrints[deposit.ContainerId], cardImage)
	}

	boxDeposits := make([]ContainerDeposit, len(boxes))

	for i, box := range boxes {
		prints, ok := boxPrints[box.ContainerId]
		if !ok {
			return nil, fmt.Errorf("unknown box %d", box.ContainerId)
		}

		total := 0
		for _, card := range prints {
			total += card.Amount
		}

		boxDeposits[i] = ContainerDeposit{
			box.ContainerId,
			box.Name,
			total,
			prints,
		}
	}

	boxesByOrder := make(map[int]int, len(boxes))
	for _, box := range boxes {
		boxesByOrder[box.ContainerId] = box.SortOrder
	}

	slices.SortFunc(boxDeposits, func(a, b ContainerDeposit) int {
		return cmp.Compare(boxesByOrder[a.ContainerId], boxesByOrder[b.ContainerId])
	})

	return boxDeposits, nil
}

func FilterCards(fullCards []cards.Card, deposits []CardDeposit) []cards.Card {
	oracleIds := map[uuid.UUID]bool{}
	for _, deposit := range deposits {
		oracleIds[deposit.OracleId] = true
	}

	matches := make([]cards.Card, len(oracleIds))
	i := 0
	for _, card := range fullCards {
		if oracleIds[card.OracleId] {
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

func ToScryfallIds(deposits []CardDeposit) uuid.UUIDs {
	uniqIds := make(map[uuid.UUID]struct{})
	var v struct{}
	for _, deposit := range deposits {
		uniqIds[deposit.ScryfallId] = v
	}

	ids := make(uuid.UUIDs, len(uniqIds))
	i := 0
	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
}

func ToContainerIds(deposits []CardDeposit) []int {
	uniqIds := make(map[int]struct{})
	var v struct{}
	for _, deposit := range deposits {
		uniqIds[deposit.ContainerId] = v
	}

	ids := make([]int, len(uniqIds))
	i := 0
	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
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
