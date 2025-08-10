package containers

import (
	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type Container struct {
	Name            string
	Capacity        int
	MarkForDeletion bool
}

type CardDeposit struct {
	ScryfallId uuid.UUID
	Amount     int
}

type CardAmount struct {
	cards.Card
	Amount int `json:"amount"`
}
