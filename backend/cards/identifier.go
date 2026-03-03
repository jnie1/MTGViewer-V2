package cards

import "github.com/google/uuid"

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
	ScryfallIdentifier | MultiverseIdentifier | SetCollectorNumber | NameSet
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
