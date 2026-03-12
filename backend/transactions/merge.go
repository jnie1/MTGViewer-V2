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
		newChange := containerChange{key.containerId, delta}
		cardId := key.scryfallId
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

		// sort for largest adds first
		slices.SortFunc(adds, func(a, b containerChange) int {
			return -cmp.Compare(a.delta, b.delta)
		})
		// sort for largest deletes first
		slices.SortFunc(deletes, func(a, b containerChange) int {
			return cmp.Compare(a.delta, b.delta)
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
				FromContainer: containers[currentDelete.containerId],
				ToContainer:   containers[currentAdd.containerId],
				ScryfallId:    cardId,
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
				FromContainer: containers[currentDelete.containerId],
				ScryfallId:    cardId,
				Quantity:      -currentDelete.delta,
			})
		}

		if j < len(deletes) {
			for _, extra := range deletes[j:] {
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

		if i < len(adds) {
			for _, extra := range adds[i:] {
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
