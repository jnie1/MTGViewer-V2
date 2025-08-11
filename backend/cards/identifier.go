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
