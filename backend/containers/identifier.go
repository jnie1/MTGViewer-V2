package containers

import (
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type ScryfallIdentifier struct {
	Id uuid.UUID `json:"scryfallId"`
}

type MultiverseIdentifier struct {
	MultiverseId int `json:"multiverseId"`
}

type SetCollectorNumber struct {
	Set             string `json:"set"`
	CollectorNumber string `json:"collectorNumber"`
}

type NameSet struct {
	Name string `json:"name"`
	Set  string `json:"set"`
}

type CardIdentifier interface {
	Copy() cards.CardIdentifier
}

func (id ScryfallIdentifier) Copy() cards.CardIdentifier {
	return cards.ScryfallIdentifier{Id: id.Id}
}

func (id MultiverseIdentifier) Copy() cards.CardIdentifier {
	return cards.MultiverseIdentifier{MultiverseId: id.MultiverseId}
}

func (id SetCollectorNumber) Copy() cards.CardIdentifier {
	return cards.SetCollectorNumber{Set: id.Set, CollectorNumber: id.CollectorNumber}
}

func (id NameSet) Copy() cards.CardIdentifier {
	return cards.NameSet{Name: id.Name, Set: id.Set}
}
