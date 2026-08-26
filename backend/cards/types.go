package cards

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type CardImageURLs struct {
	Preview string `json:"preview,omitempty"`
	Normal  string `json:"normal,omitempty"`
	Full    string `json:"full,omitempty"`
}

type Card struct {
	ScryfallId      uuid.UUID     `json:"scryfallId"`
	OracleId        uuid.UUID     `json:"oracleId"`
	Name            string        `json:"name"`
	ManaCost        string        `json:"manaCost,omitempty"`
	Set             string        `json:"set"`
	SetCode         string        `json:"setCode"`
	CollectorNumber string        `json:"collectorNumber"`
	MultiverseId    int           `json:"multiverseId,omitempty"`
	Type            string        `json:"type"`
	Rarity          string        `json:"rarity"`
	Power           string        `json:"power,omitempty"`
	Toughness       string        `json:"toughness,omitempty"`
	Images          CardImageURLs `json:"imageUrls"`
}

type CardAmount struct {
	Card
	Amount int `json:"amount"`
}

type CardAmountPreview struct {
	ScryfallId uuid.UUID `json:"scryfallId"`
	OracleId   uuid.UUID `json:"oracleId"`
	Amount     int       `json:"amount"`
}

type CardImagePreview struct {
	ScryfallId uuid.UUID     `json:"scryfallId"`
	Name       string        `json:"name"`
	Images     CardImageURLs `json:"imageUrls"`
}

type CardImageAmount struct {
	ScryfallId uuid.UUID     `json:"scryfallId"`
	Images     CardImageURLs `json:"imageUrls"`
	Amount     int           `json:"amount"`
}

type CardPricePreview struct {
	ScryfallId uuid.UUID `json:"scryfallId"`
	Price      float64   `json:"price"`
}

type CardPriceAmount struct {
	CardAmount
	Price float64 `json:"price"`
}

type SearchCardPage struct {
	Cards   []Card `json:"cards"`
	Page    int    `json:"page"`
	HasMore bool   `json:"hasMore"`
}

func ParseScryfallIds(ids []string) (uuid.UUIDs, error) {
	scryfallIds := make(uuid.UUIDs, len(ids))

	for i, id := range ids {
		id, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		scryfallIds[i] = id
	}

	return scryfallIds, nil
}

func ToScryfallIds(amounts []CardAmountPreview) uuid.UUIDs {
	uniqIds := map[uuid.UUID]bool{}

	for _, deposit := range amounts {
		uniqIds[deposit.ScryfallId] = true
	}

	ids := make(uuid.UUIDs, len(uniqIds))
	i := 0

	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
}

func JoinCardAmounts(cards []Card, previews []CardAmountPreview) ([]CardAmount, error) {
	amounts := make([]CardAmount, len(previews))
	cardMap := make(map[uuid.UUID]Card, len(cards))

	for _, card := range cards {
		cardMap[card.ScryfallId] = card
	}

	for i, deposit := range previews {
		card, ok := cardMap[deposit.ScryfallId]
		if !ok {
			return nil, fmt.Errorf("cannot resolve card id %s", deposit.ScryfallId)
		}
		amounts[i] = CardAmount{
			Card:   card,
			Amount: deposit.Amount,
		}
	}

	slices.SortFunc(amounts, func(a, b CardAmount) int {
		nameCompare := strings.Compare(a.Name, b.Name)
		if nameCompare == 0 {
			return strings.Compare(a.Set, b.Set)
		}
		return nameCompare
	})

	return amounts, nil
}
