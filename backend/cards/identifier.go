package cards

import "github.com/google/uuid"

type ScryfallIdentifier struct {
	Id uuid.UUID `json:"id"`
}

type SetCollectorNumber struct {
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
}

type CollectionQuery[Id ScryfallIdentifier | SetCollectorNumber] struct {
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
