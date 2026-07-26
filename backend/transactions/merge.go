package transactions

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/containers"
)

func MergeLogs(logs []CardLogPreview) []CardLogPreview {
	containerDeltas := map[containers.ContainerCard]int{}
	containersById := containerMappings{}

	for _, log := range logs {
		if log.FromContainer != nil {
			containerId := log.FromContainer.ContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}

			containerDeltas[key] = containerDeltas[key] - log.Amount
			containersById[containerId] = log.FromContainer
		}
		if log.ToContainer != nil {
			containerId := log.ToContainer.ContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}

			containerDeltas[key] = containerDeltas[key] + log.Amount
			containersById[containerId] = log.ToContainer
		}
	}

	changesByCard := map[uuid.UUID][]containerChange{}

	for key, delta := range containerDeltas {
		cardId := key.ScryfallId
		newChange := containerChange{key.ContainerId, delta}
		changesByCard[cardId] = append(changesByCard[cardId], newChange)
	}

	return combineCardDeltas(changesByCard, containersById)
}

func combineCardDeltas(deltas map[uuid.UUID][]containerChange, containers containerMappings) []CardLogPreview {
	updatedLogs := []CardLogPreview{}

	for cardId, changes := range deltas {
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
			newLog := CardLogPreview{
				FromContainer: containers[currentDelete.containerId],
				ToContainer:   containers[currentAdd.containerId],
				ScryfallId:    cardId,
			}

			if add < delete {
				newLog.Amount = add

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
				newLog.Amount = delete

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
				newLog.Amount = add

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

		if currentAdd.delta > 0 {
			updatedLogs = append(updatedLogs, CardLogPreview{
				ToContainer: containers[currentAdd.containerId],
				ScryfallId:  cardId,
				Amount:      currentAdd.delta,
			})

			if addIndex+1 < len(adds) {
				for _, extra := range adds[addIndex+1:] {
					updatedLogs = append(updatedLogs, CardLogPreview{
						ToContainer: containers[extra.containerId],
						ScryfallId:  cardId,
						Amount:      extra.delta,
					})
				}
			}
		}

		if currentDelete.delta < 0 {
			updatedLogs = append(updatedLogs, CardLogPreview{
				FromContainer: containers[currentDelete.containerId],
				ScryfallId:    cardId,
				Amount:        -currentDelete.delta,
			})

			if deleteIndex+1 < len(deletes) {
				for _, extra := range deletes[deleteIndex+1:] {
					updatedLogs = append(updatedLogs, CardLogPreview{
						FromContainer: containers[extra.containerId],
						ScryfallId:    cardId,
						Amount:        -extra.delta,
					})
				}
			}
		}
	}

	return updatedLogs
}
