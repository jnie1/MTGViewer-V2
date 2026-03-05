package containers

import (
	"errors"
	"math"

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

var ErrUnknownCardIdentifier = errors.New("unknown card identifier specified")

func FromObj(obj map[string]any) (CardIdentifier, error) {
	if str, ok := obj["scryfallId"].(string); ok {
		scryfallId, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		if len(obj) == 1 {
			return ScryfallIdentifier{scryfallId}, nil
		}
	}

	if multiverseId, ok := obj["multiverseId"].(float64); ok {
		if multiverseId == math.Trunc(multiverseId) {
			if len(obj) == 1 {
				return MultiverseIdentifier{int(multiverseId)}, nil
			}
		}
	}

	if collectorNumber, ok := obj["collectorNumber"].(string); ok {
		if set, ok := obj["set"].(string); ok {
			if len(obj) == 2 {
				return SetCollectorNumber{set, collectorNumber}, nil
			}
		}
	}

	if name, ok := obj["name"].(string); ok {
		if set, ok := obj["set"].(string); ok {
			if len(obj) == 2 {
				return NameSet{name, set}, nil
			}
		}
	}

	return nil, ErrUnknownCardIdentifier
}
