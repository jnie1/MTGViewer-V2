package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type LogRange struct {
	Start time.Time
	End   time.Time
}

type CardTransaction struct {
	GroupId uuid.UUID `json:"groupId"`
	Time    time.Time `json:"time"`
	Total   int       `json:"total"`
}

type CardLogPreview struct {
	ScryfallId      uuid.UUID
	FromContainerId *int
	ToContainerId   *int
	Amount          int
}

type CardTransfer struct {
	cards.CardImagePreview
	Delta           int  `json:"delta"`
	WithContainerId *int `json:"withContainerId,omitempty"`
}

type ContainerTransfers struct {
	ContainerId   int            `json:"containerId"`
	ContainerName string         `json:"containerName"`
	Total         int            `json:"total"`
	Cards         []CardTransfer `json:"cards"`
}

func ToScryfallIds(transactionLogs []CardLogPreview) uuid.UUIDs {
	uniqIds := make(map[uuid.UUID]struct{})
	var v struct{}
	for _, log := range transactionLogs {
		uniqIds[log.ScryfallId] = v
	}

	ids := make(uuid.UUIDs, len(uniqIds))
	i := 0
	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
}

func ToContainerIds(transactionLogs []CardLogPreview) []int {
	uniqIds := make(map[int]struct{})
	var v struct{}
	for _, log := range transactionLogs {
		if log.FromContainerId != nil {
			uniqIds[*log.FromContainerId] = v
		}
		if log.ToContainerId != nil {
			uniqIds[*log.ToContainerId] = v
		}
	}

	ids := make([]int, len(uniqIds))
	i := 0
	for id := range uniqIds {
		ids[i] = id
		i += 1
	}

	return ids
}
