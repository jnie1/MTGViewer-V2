package transactions

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/containers"
)

func MergeLogs(logs []CardLogPreview) []CardLogPreview {
	changesByContainer := map[containers.ContainerCard]int{}
	containersById := map[int]*containers.ContainerPreview{}

	for _, log := range logs {
		if log.FromContainer != nil {
			containerId := log.FromContainer.ContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}

			changesByContainer[key] = changesByContainer[key] - log.Amount
			containersById[containerId] = log.FromContainer
		}
		if log.ToContainer != nil {
			containerId := log.ToContainer.ContainerId
			key := containers.ContainerCard{ContainerId: containerId, ScryfallId: log.ScryfallId}

			changesByContainer[key] = changesByContainer[key] + log.Amount
			containersById[containerId] = log.ToContainer
		}
	}

	changesByCard := map[uuid.UUID][]containers.ContainerDelta{}
	for cardKey, delta := range changesByContainer {
		cardId := cardKey.ScryfallId
		newChange := containers.ContainerDelta{ContainerId: cardKey.ContainerId, Delta: delta}
		changesByCard[cardId] = append(changesByCard[cardId], newChange)
	}

	updatedLogs := []CardLogPreview{}

	for cardId, changes := range changesByCard {
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

		var currentAdd, currentDelete containers.ContainerDelta
		addIndex, deleteIndex := 0, 0

		if addIndex < len(adds) {
			currentAdd = adds[addIndex]
		}
		if deleteIndex < len(dels) {
			currentDelete = dels[deleteIndex]
		}

		for currentAdd.Delta > 0 && currentDelete.Delta < 0 {
			add, del := currentAdd.Delta, -currentDelete.Delta
			newLog := CardLogPreview{
				FromContainer: containersById[currentDelete.ContainerId],
				ToContainer:   containersById[currentAdd.ContainerId],
				ScryfallId:    cardId,
			}

			if add < del {
				newLog.Amount = add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containers.ContainerDelta{}
				}

				currentDelete = containers.ContainerDelta{
					ContainerId: currentDelete.ContainerId,
					Delta:       currentDelete.Delta + add,
				}
			} else if add > del {
				newLog.Amount = del

				deleteIndex += 1
				if deleteIndex < len(dels) {
					currentDelete = dels[deleteIndex]
				} else {
					currentDelete = containers.ContainerDelta{}
				}

				currentAdd = containers.ContainerDelta{
					ContainerId: currentAdd.ContainerId,
					Delta:       currentAdd.Delta - del,
				}
			} else {
				newLog.Amount = add

				addIndex += 1
				if addIndex < len(adds) {
					currentAdd = adds[addIndex]
				} else {
					currentAdd = containers.ContainerDelta{}
				}

				deleteIndex += 1
				if deleteIndex < len(dels) {
					currentDelete = dels[deleteIndex]
				} else {
					currentDelete = containers.ContainerDelta{}
				}
			}

			updatedLogs = append(updatedLogs, newLog)
		}

		if currentAdd.Delta > 0 {
			updatedLogs = append(updatedLogs, CardLogPreview{
				ToContainer: containersById[currentAdd.ContainerId],
				ScryfallId:  cardId,
				Amount:      currentAdd.Delta,
			})

			if addIndex+1 < len(adds) {
				for _, extra := range adds[addIndex+1:] {
					updatedLogs = append(updatedLogs, CardLogPreview{
						ToContainer: containersById[extra.ContainerId],
						ScryfallId:  cardId,
						Amount:      extra.Delta,
					})
				}
			}
		}

		if currentDelete.Delta < 0 {
			updatedLogs = append(updatedLogs, CardLogPreview{
				FromContainer: containersById[currentDelete.ContainerId],
				ScryfallId:    cardId,
				Amount:        -currentDelete.Delta,
			})

			if deleteIndex+1 < len(dels) {
				for _, extra := range dels[deleteIndex+1:] {
					updatedLogs = append(updatedLogs, CardLogPreview{
						FromContainer: containersById[extra.ContainerId],
						ScryfallId:    cardId,
						Amount:        -extra.Delta,
					})
				}
			}
		}
	}

	return updatedLogs
}
