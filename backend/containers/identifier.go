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

func (si ScryfallIdentifier) Copy() cards.CardIdentifier {
	return cards.ScryfallIdentifier{Id: si.Id}
}

func (mi MultiverseIdentifier) Copy() cards.CardIdentifier {
	return cards.MultiverseIdentifier{MultiverseId: mi.MultiverseId}
}

func (sc SetCollectorNumber) Copy() cards.CardIdentifier {
	return cards.SetCollectorNumber{Set: sc.Set, CollectorNumber: sc.CollectorNumber}
}

func (ns NameSet) Copy() cards.CardIdentifier {
	return cards.NameSet{Name: ns.Name, Set: ns.Set}
}
