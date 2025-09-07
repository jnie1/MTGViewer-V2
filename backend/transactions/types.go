package transactions

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

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

func (container *TransactionContainer) GetContainer() TransactionContainer {
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
	cardMap := map[uuid.UUID]cards.Card{}

	for _, loggedCard := range loggedCards {
		cardMap[loggedCard.ScryfallId] = loggedCard
	}

	reportCards := make([]ReportCard, len(logs))

	for i, log := range logs {
		if reportedCard, ok := cardMap[log.ScryfallId]; ok {
			reportCards[i] = ReportCard{
				GroupId:       log.GroupId,
				FromContainer: log.FromContainer,
				ToContainer:   log.ToContainer,
				Card:          reportedCard,
				Quantity:      log.Quantity,
			}
		} else {
			return nil, fmt.Errorf("cannot resolve card id %s", log.ScryfallId)
		}
	}

	slices.SortFunc(reportCards, compareReportCards)

	return reportCards, nil
}

func compareReportCards(a, b ReportCard) int {
	if c := cmp.Compare(a.FromContainer.GetContainer().Name, b.FromContainer.GetContainer().Name); c != 0 {
		return c
	}

	if c := cmp.Compare(a.ToContainer.GetContainer().Name, b.ToContainer.GetContainer().Name); c != 0 {
		return c
	}

	return cmp.Compare(a.Card.Name, b.Card.Name)
}
