package transactions

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
)

func MergeLogs(logs []TransactionLogs) []TransactionLogs {
	containerDeltas := map[containerCard]int{}
	containersById := map[int]*TransactionContainer{}

	for _, log := range logs {
		if log.FromContainer != nil {
			containerId := log.FromContainer.ContainerId
			key := containerCard{containerId, log.ScryfallId}

			containerDeltas[key] = containerDeltas[key] - log.Quantity
			containersById[containerId] = log.FromContainer
		}
		if log.ToContainer != nil {
			containerId := log.ToContainer.ContainerId
			key := containerCard{containerId, log.ScryfallId}

			containerDeltas[key] = containerDeltas[key] + log.Quantity
			containersById[containerId] = log.ToContainer
		}
	}

	if len(containerDeltas) == len(logs) {
		return logs
	}

	return combineContainerDeltas(containerDeltas, containersById)
}

func combineContainerDeltas(deltas map[containerCard]int, containers map[int]*TransactionContainer) []TransactionLogs {
	updatedLogs := make([]TransactionLogs, len(deltas))
	changesByCard := map[uuid.UUID][]containerChange{}

	for key, delta := range deltas {
		cardId := key.scryfallId
		newChange := containerChange{key.containerId, delta}
		changesByCard[cardId] = append(changesByCard[cardId], newChange)
	}

	for cardId, changes := range changesByCard {
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
			// sort largest adds first, the most positive number desc
			return -cmp.Compare(a.delta, b.delta)
		})
		slices.SortFunc(deletes, func(a, b containerChange) int {
			// sort largest deletes first, most negative number asc
			return cmp.Compare(a.delta, b.delta)
		})

		var currentAdd, currentDelete containerChange
		addIndex, deleteIndex := 0, 0

		if addIndex < len(adds) {
			currentAdd = adds[addIndex]
		}
		if deleteIndex < len(deletes) {
			currentDelete = deletes[deleteIndex]
		}

		for currentAdd.delta > 0 && currentDelete.delta < 0 {
			add, delete := currentAdd.delta, -currentDelete.delta
			newLog := TransactionLogs{
				FromContainer: containers[currentDelete.containerId],
				ToContainer:   containers[currentAdd.containerId],
				ScryfallId:    cardId,
			}

			if add < delete {
				newLog.Quantity = add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containerChange{}
				}

				currentDelete = containerChange{
					containerId: currentDelete.containerId,
					delta:       currentDelete.delta + add,
				}
			} else if add > delete {
				newLog.Quantity = delete

				deleteIndex += 1
				if deleteIndex < len(deletes) {
					currentDelete = deletes[deleteIndex]
				} else {
					currentDelete = containerChange{}
				}

				currentAdd = containerChange{
					containerId: currentAdd.containerId,
					delta:       currentAdd.delta - delete,
				}
			} else {
				newLog.Quantity = add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containerChange{}
				}

				deleteIndex += 1
				if deleteIndex < len(deletes) {
					currentDelete = deletes[deleteIndex]
				} else {
					currentDelete = containerChange{}
				}
			}

			updatedLogs = append(updatedLogs, newLog)
		}

		if currentDelete.delta < 0 {
			updatedLogs = append(updatedLogs, TransactionLogs{
				FromContainer: containers[currentDelete.containerId],
				ScryfallId:    cardId,
				Quantity:      -currentDelete.delta,
			})
		}

		if deleteIndex < len(deletes) {
			for _, extra := range deletes[deleteIndex:] {
				updatedLogs = append(updatedLogs, TransactionLogs{
					FromContainer: containers[extra.containerId],
					ScryfallId:    cardId,
					Quantity:      -extra.delta,
				})
			}
		}

		if currentAdd.delta > 0 {
			updatedLogs = append(updatedLogs, TransactionLogs{
				ToContainer: containers[currentAdd.containerId],
				ScryfallId:  cardId,
				Quantity:    currentAdd.delta,
			})
		}

		if addIndex < len(adds) {
			for _, extra := range adds[addIndex:] {
				updatedLogs = append(updatedLogs, TransactionLogs{
					ToContainer: containers[extra.containerId],
					ScryfallId:  cardId,
					Quantity:    extra.delta,
				})
			}
		}
	}

	return updatedLogs
}
