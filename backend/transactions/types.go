package transactions

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

type UpdateLogs struct {
	GroupId uuid.UUID `json:"groupId"`
	Time    time.Time `json:"time"`
}

type TransactionLogs struct {
	TransactionId int                   `json:"transactionId"`
	GroupId       uuid.UUID             `json:"groupId"`
	FromContainer *TransactionContainer `json:"fromContainer"`
	ToContainer   *TransactionContainer `json:"toContainer"`
	ScryfallId    uuid.UUID             `json:"scryfallId"`
	Quantity      int                   `json:"quantity"`
}

type TransactionContainer struct {
	ContainerId int    `json:"containerId"`
	Name        string `json:"name"`
}

func (container *TransactionContainer) Container() TransactionContainer {
	if container == nil {
		return TransactionContainer{}
	}
	return *container
}

type ReportCard struct {
	GroupId       uuid.UUID             `json:"groupId"`
	FromContainer *TransactionContainer `json:"fromContainer"`
	ToContainer   *TransactionContainer `json:"toContainer"`
	Card          cards.Card            `json:"card"`
	Quantity      int                   `json:"quantity"`
}

func GetScryfallIds(transactionLogs []TransactionLogs) []cards.ScryfallIdentifier {
	uniqIds := map[uuid.UUID]any{}

	for _, log := range transactionLogs {
		uniqIds[log.ScryfallId] = nil
	}

	allIds := make([]cards.ScryfallIdentifier, len(uniqIds))
	i := 0

	for id := range uniqIds {
		allIds[i] = cards.ScryfallIdentifier{Id: id}
		i += 1
	}

	return allIds
}

func JoinReportCards(loggedCards []cards.Card, logs []TransactionLogs) ([]ReportCard, error) {
	reportCards := make([]ReportCard, len(logs))
	cardMap := make(map[uuid.UUID]cards.Card, len(loggedCards))

	for _, loggedCard := range loggedCards {
		cardMap[loggedCard.ScryfallId] = loggedCard
	}

	for i, log := range logs {
		reportedCard, ok := cardMap[log.ScryfallId]
		if !ok {
			return nil, fmt.Errorf("cannot resolve card id %s", log.ScryfallId)
		}
		reportCards[i] = ReportCard{
			GroupId:       log.GroupId,
			FromContainer: log.FromContainer,
			ToContainer:   log.ToContainer,
			Card:          reportedCard,
			Quantity:      log.Quantity,
		}
	}

	slices.SortFunc(reportCards, compareReportCards)
	return reportCards, nil
}

func compareReportCards(a, b ReportCard) int {
	if c := cmp.Compare(a.FromContainer.Container().Name, b.FromContainer.Container().Name); c != 0 {
		return c
	}

	if c := cmp.Compare(a.ToContainer.Container().Name, b.ToContainer.Container().Name); c != 0 {
		return c
	}

	return cmp.Compare(a.Card.Name, b.Card.Name)
}
