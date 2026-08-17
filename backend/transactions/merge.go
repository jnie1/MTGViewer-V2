package transactions

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
	"github.com/jnie1/MTGViewer-V2/containers"
)

func MergeLogs(logs []CardLogPreview, boxes []containers.ContainerPreview, fullCards []cards.Card) ([]ContainerTransfers, error) {
	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}

	transfersByBox, err := combineTransfersByBox(logs, boxes, cardsById)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(boxes, func(a, b containers.ContainerPreview) int {
		return cmp.Compare(a.SortOrder, b.SortOrder)
	})

	boxSorts := make(map[int]int, len(boxes))
	for _, box := range boxes {
		boxSorts[box.ContainerId] = box.SortOrder
	}

	transfers := make([]ContainerTransfers, len(boxes))

	for i, box := range boxes {
		logs, ok := transfersByBox[box.ContainerId]
		if !ok {
			return nil, fmt.Errorf("unknown box %d", box.ContainerId)
		}

		total := 0
		for _, transfer := range logs {
			total += transfer.Delta
		}

		slices.SortFunc(logs, func(a, b CardTransfer) int {
			cardA := cardsById[a.ScryfallId]
			cardB := cardsById[b.ScryfallId]

			cardCmp := cmp.Compare(cardA.Name, cardB.Name)
			if cardCmp != 0 {
				return cardCmp
			}

			var boxA, boxB int
			if a.WithContainerId != nil {
				boxA = boxSorts[*a.WithContainerId]
			}
			if b.WithContainerId != nil {
				boxB = boxSorts[*b.WithContainerId]
			}

			return cmp.Compare(boxA, boxB)
		})

		transfers[i] = ContainerTransfers{
			ContainerId:   box.ContainerId,
			ContainerName: box.Name,
			Total:         total,
			Cards:         logs,
		}
	}

	return transfers, nil
}

func groupDeltasByCard(logs []CardLogPreview) map[uuid.UUID][]containers.ContainerDelta {
	changesByContainer := map[containers.ContainerCard]int{}

	for _, log := range logs {
		if log.FromContainerId != nil {
			containerId := *log.FromContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}
			changesByContainer[key] = changesByContainer[key] - log.Amount
		}
		if log.ToContainerId != nil {
			containerId := *log.ToContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}
			changesByContainer[key] = changesByContainer[key] + log.Amount
		}
	}

	changesByCard := map[uuid.UUID][]containers.ContainerDelta{}

	for cardKey, delta := range changesByContainer {
		cardId := cardKey.ScryfallId
		newChange := containers.ContainerDelta{ContainerId: cardKey.ContainerId, Delta: delta}
		changesByCard[cardId] = append(changesByCard[cardId], newChange)
	}

	return changesByCard
}

func combineTransfersByBox(logs []CardLogPreview, boxes []containers.ContainerPreview, cardsById map[uuid.UUID]cards.Card) (map[int][]CardTransfer, error) {
	transfersByBox := make(map[int][]CardTransfer, len(boxes))

	for cardId, changes := range groupDeltasByCard(logs) {
		fullCard, ok := cardsById[cardId]
		if !ok {
			return nil, fmt.Errorf("unknown card %s", cardId)
		}

		card := cards.CardImagePreview{
			ScryfallId: fullCard.ScryfallId,
			Name:       fullCard.Name,
			Images:     fullCard.Images,
		}

		adds := []containers.ContainerDelta{}
		dels := []containers.ContainerDelta{}

		for _, change := range changes {
			if change.Delta > 0 {
				adds = append(adds, change)
			} else if change.Delta < 0 {
				dels = append(dels, change)
			}
		}

		slices.SortFunc(adds, func(a, b containers.ContainerDelta) int {
			// sort largest adds first, the most positive number desc
			return -cmp.Compare(a.Delta, b.Delta)
		})
		slices.SortFunc(dels, func(a, b containers.ContainerDelta) int {
			// sort largest deletes first, most negative number asc
			return cmp.Compare(a.Delta, b.Delta)
		})

		var currentAdd, currentDel containers.ContainerDelta
		addIndex, deleteIndex := 0, 0

		if addIndex < len(adds) {
			currentAdd = adds[addIndex]
		}
		if deleteIndex < len(dels) {
			currentDel = dels[deleteIndex]
		}

		for currentAdd.Delta > 0 && currentDel.Delta < 0 {
			add, del := currentAdd.Delta, -currentDel.Delta

			addContainer := currentAdd.ContainerId
			addTransfer := CardTransfer{card, 0, &currentDel.ContainerId}

			delContainer := currentDel.ContainerId
			delTransfer := CardTransfer{card, 0, &currentAdd.ContainerId}

			if add < del {
				addTransfer.Delta = add
				delTransfer.Delta = -add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containers.ContainerDelta{}
				}

				currentDel = containers.ContainerDelta{
					ContainerId: currentDel.ContainerId,
					Delta:       currentDel.Delta + add,
				}
			} else if add > del {
				addTransfer.Delta = del
				delTransfer.Delta = -del

				deleteIndex += 1
				if deleteIndex < len(dels) {
					currentDel = dels[deleteIndex]
				} else {
					currentDel = containers.ContainerDelta{}
				}

				currentAdd = containers.ContainerDelta{
					ContainerId: currentAdd.ContainerId,
					Delta:       currentAdd.Delta - del,
				}
			} else {
				addTransfer.Delta = add
				delTransfer.Delta = -add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containers.ContainerDelta{}
				}

				deleteIndex += 1
				if deleteIndex < len(dels) {
					currentDel = dels[deleteIndex]
				} else {
					currentDel = containers.ContainerDelta{}
				}
			}

			transfersByBox[addContainer] = append(transfersByBox[addContainer], addTransfer)
			transfersByBox[delContainer] = append(transfersByBox[delContainer], delTransfer)
		}

		if currentAdd.Delta > 0 {
			id := currentAdd.ContainerId
			transfersByBox[id] = append(transfersByBox[id], CardTransfer{card, currentAdd.Delta, nil})

			if addIndex+1 < len(adds) {
				for _, extra := range adds[addIndex+1:] {
					id := extra.ContainerId
					transfersByBox[id] = append(transfersByBox[id], CardTransfer{card, extra.Delta, nil})
				}
			}
		}

		if currentDel.Delta < 0 {
			id := currentDel.ContainerId
			transfersByBox[id] = append(transfersByBox[id], CardTransfer{card, currentDel.Delta, nil})

			if deleteIndex+1 < len(dels) {
				for _, extra := range dels[deleteIndex+1:] {
					id := extra.ContainerId
					transfersByBox[id] = append(transfersByBox[id], CardTransfer{card, extra.Delta, nil})
				}
			}
		}
	}

	return transfersByBox, nil
}
