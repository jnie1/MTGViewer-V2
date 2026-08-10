package cards

import (
	"errors"
	"math"

	"github.com/google/uuid"
)

type ScryfallIdObj struct {
	ScryfallId uuid.UUID `json:"scryfallId"`
}

type MultiverseIdObj struct {
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

type CardId struct {
	ScryfallId      uuid.UUID `json:"scryfallId"`
	Name            string    `json:"name" `
	SetCode         string    `json:"setCode"`
	CollectorNumber string    `json:"collectorNumber"`
	MultiverseId    int       `json:"multiverseId,omitempty"`
}

func (id CardId) SetCollectorNumber() SetCollectorNumber {
	return SetCollectorNumber{id.SetCode, id.CollectorNumber}
}

func (id CardId) NameSet() NameSet {
	return NameSet{id.Name, id.SetCode}
}

var ErrUnknownCardIdentifier = errors.New("unknown card identifier specified")

func FromObj(obj map[string]any) (any, error) {
	if str, ok := obj["scryfallId"].(string); ok {
		scryfallId, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		if len(obj) == 1 {
			return ScryfallIdObj{scryfallId}, nil
		}
	}

	if multiverseId, ok := obj["multiverseId"].(float64); ok {
		if multiverseId == math.Trunc(multiverseId) {
			if len(obj) == 1 {
				return MultiverseIdObj{int(multiverseId)}, nil
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
