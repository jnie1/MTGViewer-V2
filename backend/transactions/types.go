package transactions

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
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
	FromContainer *containers.ContainerName
	ToContainer   *containers.ContainerName
	ScryfallId    uuid.UUID
	Amount        int
}

type CardLog struct {
	FromContainer *containers.ContainerName `json:"fromContainer"`
	ToContainer   *containers.ContainerName `json:"toContainer"`
	Card          cards.Card                `json:"card"`
	Amount        int                       `json:"amount"`
}

type containerChange struct {
	containerId int
	delta       int
}

type containerMappings map[int]*containers.ContainerName

func ToScryfallIds(transactionLogs []CardLogPreview) []cards.ScryfallIdentifier {
	uniqIds := map[uuid.UUID]any{}

	for _, log := range transactionLogs {
		uniqIds[log.ScryfallId] = nil
	}

	ids := make([]cards.ScryfallIdentifier, len(uniqIds))
	i := 0

	for id := range uniqIds {
		ids[i] = cards.ScryfallIdentifier{Id: id}
		i += 1
	}

	return ids
}

func JoinCardLogs(loggedCards []cards.Card, logs []CardLogPreview) ([]CardLog, error) {
	cardChanges := make([]CardLog, len(logs))
	cardMap := make(map[uuid.UUID]cards.Card, len(loggedCards))

	for _, loggedCard := range loggedCards {
		cardMap[loggedCard.ScryfallId] = loggedCard
	}

	for i, log := range logs {
		reportedCard, ok := cardMap[log.ScryfallId]
		if !ok {
			return nil, fmt.Errorf("cannot resolve card id %s", log.ScryfallId)
		}
		cardChanges[i] = CardLog{
			FromContainer: log.FromContainer,
			ToContainer:   log.ToContainer,
			Card:          reportedCard,
			Amount:        log.Amount,
		}
	}

	slices.SortFunc(cardChanges, compareCardChange)
	return cardChanges, nil
}

func compareCardChange(a, b CardLog) int {
	if c := cmp.Compare(a.FromContainer.Container().Name, b.FromContainer.Container().Name); c != 0 {
		return c
	}

	if c := cmp.Compare(a.ToContainer.Container().Name, b.ToContainer.Container().Name); c != 0 {
		return c
	}

	return cmp.Compare(a.Card.Name, b.Card.Name)
}
