package containers

import (
	"github.com/google/uuid"
)

type Container struct {
	Name            string
	Capacity        int
	MarkForDeletion bool
}

type CardDeposit struct {
	ContainerId int
	ScryfallId  uuid.UUID
	Amount      int
}

type CardRequest struct {
	ScryfallId uuid.UUID
	Delta      int
}

type ContainerAllocation struct {
	ContainerId int
	Used        int
	MaxCapcity  int
}

type ContainerChanges struct {
	ContainerId int
	Requests    []CardRequest
}
