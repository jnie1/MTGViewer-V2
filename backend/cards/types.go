package cards

import (
	"strings"

	"github.com/google/uuid"
)

type CardImageUrls struct {
	Preview string `json:"preview,omitempty"`
	Normal  string `json:"normal,omitempty"`
	Full    string `json:"full,omitempty"`
}

type Card struct {
	ScryfallId      uuid.UUID     `json:"scryfallId"`
	Name            string        `json:"name"`
	ManaCost        string        `json:"manaCost,omitempty"`
	Set             string        `json:"set"`
	SetCode         string        `json:"set_code"`
	CollectorNumber string        `json:"collector_number"`
	MultiverseIds   []int         `json:"multiverse_ids,omitempty"`
	Type            string        `json:"type"`
	Rarity          string        `json:"rarity"`
	Power           string        `json:"power,omitempty"`
	Toughness       string        `json:"toughness,omitempty"`
	Images          CardImageUrls `json:"imageUrls"`
}

type CardAmount struct {
	Card
	Amount int `json:"amount"`
}

type SearchCardPage struct {
	TotalCards int
	Cards      []Card
	Page       int
	HasMore    bool
}

type scryfallImages struct {
	Small  string `json:"small,omitempty"`
	Normal string `json:"normal,omitempty"`
	Large  string `json:"large,omitempty"`
}

type scryfallCardFace struct {
	Name     string         `json:"name"`
	ManaCost string         `json:"mana_cost,omitempty"`
	Type     string         `json:"type_line"`
	Images   scryfallImages `json:"image_uris"`
}

type scryfallCard struct {
	ScryfallId      uuid.UUID          `json:"id"`
	ManaCost        string             `json:"mana_cost,omitempty"`
	Name            string             `json:"name"`
	SetName         string             `json:"set_name"`
	Set             string             `json:"set"`
	CollectorNumber string             `json:"collector_number"`
	MultiverseIds   []int              `json:"multiverse_ids,omitempty"`
	Power           string             `json:"power,omitempty"`
	Toughness       string             `json:"toughness,omitempty"`
	Images          scryfallImages     `json:"image_uris"`
	CardFaces       []scryfallCardFace `json:"card_faces,omitempty"`
	Type            string             `json:"type_line"`
	Rarity          string             `json:"rarity"`
}

type searchResult struct {
	TotalCards int            `json:"total_cards"`
	HasMore    bool           `json:"has_more"`
	Cards      []scryfallCard `json:"data"`
}

type collectionResult struct {
	Cards []scryfallCard `json:"data"`
}

type collectionBatchResult struct {
	Cards []Card
	Error error
}

func toCard(card scryfallCard) Card {
	images := card.Images
	if card.CardFaces != nil && len(card.CardFaces) > 0 {
		card.Images = card.CardFaces[0].Images
	}
	return Card{
		card.ScryfallId,
		card.Name,
		card.ManaCost,
		card.SetName,
		strings.ToUpper(card.Set),
		card.CollectorNumber,
		card.MultiverseIds,
		card.Type,
		card.Rarity,
		card.Power,
		card.Toughness,
		CardImageUrls{
			images.Small,
			images.Normal,
			images.Large,
		},
	}
}

func toCards(cards []scryfallCard) []Card {
	result := make([]Card, len(cards))
	for i, card := range cards {
		result[i] = toCard(card)
	}
	return result
}
