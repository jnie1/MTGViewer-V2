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
	Transfers     []CardTransfer `json:"transfers"`
}

func ToScryfallIds(transactionLogs []CardLogPreview) uuid.UUIDs {
	uniqIds := map[uuid.UUID]any{}

	for _, log := range transactionLogs {
		uniqIds[log.ScryfallId] = nil
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
	uniqIds := map[int]any{}

	for _, log := range transactionLogs {
		if log.FromContainerId != nil {
			uniqIds[*log.FromContainerId] = nil
		}
		if log.ToContainerId != nil {
			uniqIds[*log.ToContainerId] = nil
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
