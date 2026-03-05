package cards

import (
	"fmt"

	"github.com/google/uuid"
)

type ScryfallIdentifier struct {
	Id uuid.UUID `json:"id"`
}

type MultiverseIdentifier struct {
	MultiverseId int `json:"multiverse_id"`
}

type SetCollectorNumber struct {
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
}

type NameSet struct {
	Name string `json:"name"`
	Set  string `json:"set"`
}

type CardIdentifier interface {
	Convert(card Card) (CardIdentifier, error)
}

func (si ScryfallIdentifier) Convert(card Card) (CardIdentifier, error) {
	return ScryfallIdentifier{card.ScryfallId}, nil
}

func (mi MultiverseIdentifier) Convert(card Card) (CardIdentifier, error) {
	if len(card.MultiverseIds) == 0 {
		return nil, fmt.Errorf("card resolved with no multiverse id: %s, (%s) %s", card.Name, card.SetCode, card.CollectorNumber)
	}
	return MultiverseIdentifier{card.MultiverseIds[0]}, nil
}

func (sc SetCollectorNumber) Convert(card Card) (CardIdentifier, error) {
	return SetCollectorNumber{card.SetCode, card.CollectorNumber}, nil
}

func (ns NameSet) Convert(card Card) (CardIdentifier, error) {
	return NameSet{card.Name, card.SetCode}, nil
}

type CollectionQuery[Id CardIdentifier] struct {
	Identifiers []Id `json:"identifiers"`
}

func ParseScryfallIds(ids []string) ([]ScryfallIdentifier, error) {
	scryfallIds := make([]ScryfallIdentifier, len(ids))

	for i, id := range ids {
		id, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		scryfallIds[i] = ScryfallIdentifier{Id: id}
	}

	return scryfallIds, nil
}
