package transactions

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
)

type containerCard struct {
	containerId int
	scryfallId  uuid.UUID
}

type containerChange struct {
	containerId int
	delta       int
}

func MergeLogs(logs []TransactionLogs) []TransactionLogs {
	combinedDeltas := map[containerCard]int{}
	containersById := map[int]*TransactionContainer{}

	for _, log := range logs {
		if log.FromContainer != nil {
			containerId := log.FromContainer.ContainerId
			key := containerCard{containerId, log.ScryfallId}

			combinedDeltas[key] = combinedDeltas[key] - log.Quantity
			containersById[containerId] = log.FromContainer
		}
		if log.ToContainer != nil {
			containerId := log.ToContainer.ContainerId
			key := containerCard{containerId, log.ScryfallId}

			combinedDeltas[key] = combinedDeltas[key] + log.Quantity
			containersById[containerId] = log.ToContainer
		}
	}

	if len(combinedDeltas) == len(logs) {
		return logs
	}

	changesPerCard := map[uuid.UUID][]containerChange{}

	for key, delta := range combinedDeltas {
		newChange := containerChange{key.containerId, delta}
		cardId := key.scryfallId
		changesPerCard[cardId] = append(changesPerCard[cardId], newChange)
	}

	updatedLogs := make([]TransactionLogs, len(combinedDeltas))

	for cardId, changes := range changesPerCard {
		adds := []containerChange{}
		deletes := []containerChange{}

		for _, change := range changes {
			if change.delta > 0 {
				adds = append(adds, change)
			} else if change.delta < 0 {
				deletes = append(deletes, change)
			}
		}

		slices.SortFunc(adds, func(a, b containerChange) int {
			return cmp.Compare(a.delta, b.delta)
		})
		slices.SortFunc(deletes, func(a, b containerChange) int {
			return -cmp.Compare(a.delta, b.delta)
		})

		var currentAdd, currentDelete containerChange
		if len(adds) > 0 {
			currentAdd = adds[0]
		}
		if len(deletes) > 0 {
			currentDelete = deletes[0]
		}

		i, j := 0, 0
		for currentAdd.delta > 0 && currentDelete.delta < 0 {
			add, delete := currentAdd.delta, -currentDelete.delta
			newLog := TransactionLogs{
				ScryfallId:    cardId,
				FromContainer: containersById[currentDelete.containerId],
				ToContainer:   containersById[currentAdd.containerId],
			}

			if add < delete {
				newLog.Quantity = add

				i += 1
				if i < len(adds) {
					currentAdd = adds[i]
				} else {
					currentAdd = containerChange{}
				}

				currentDelete = containerChange{
					containerId: currentDelete.containerId,
					delta:       currentDelete.delta + add,
				}
			} else if add > delete {
				newLog.Quantity = delete

				j += 1
				if j < len(deletes) {
					currentDelete = deletes[j]
				} else {
					currentDelete = containerChange{}
				}

				currentAdd = containerChange{
					containerId: currentAdd.containerId,
					delta:       currentAdd.delta - delete,
				}
			} else {
				newLog.Quantity = add

				i += 1
				if i < len(adds) {
					currentAdd = adds[i]
				} else {
					currentAdd = containerChange{}
				}

				j += 1
				if j < len(deletes) {
					currentDelete = deletes[j]
				} else {
					currentDelete = containerChange{}
				}
			}

			updatedLogs = append(updatedLogs, newLog)
		}

		if currentDelete.delta < 0 {
			updatedLogs = append(updatedLogs, TransactionLogs{
				ScryfallId:    cardId,
				FromContainer: containersById[currentDelete.containerId],
				Quantity:      -currentDelete.delta,
			})
		}

		if j < len(deletes) {
			for _, extra := range deletes[j:] {
				updatedLogs = append(updatedLogs, TransactionLogs{
					ScryfallId:    cardId,
					FromContainer: containersById[extra.containerId],
					Quantity:      -extra.delta,
				})
			}
		}

		if currentAdd.delta > 0 {
			updatedLogs = append(updatedLogs, TransactionLogs{
				ScryfallId:  cardId,
				ToContainer: containersById[currentAdd.containerId],
				Quantity:    currentAdd.delta,
			})
		}

		if i < len(adds) {
			for _, extra := range adds[i:] {
				updatedLogs = append(updatedLogs, TransactionLogs{
					ScryfallId:  cardId,
					ToContainer: containersById[extra.containerId],
					Quantity:    extra.delta,
				})
			}
		}
	}

	return updatedLogs
}
