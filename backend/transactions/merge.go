package transactions

import "github.com/google/uuid"

type changeKey struct {
	scryfallId  uuid.UUID
	containerId int
}

type transactionDelta struct {
	transactionId int
	delta         int
}

func MergeLogs(logs []TransactionLogs) []TransactionLogs {
	changes := map[changeKey]transactionDelta{}

	for _, log := range logs {
		if log.FromContainer != nil {
			key := changeKey{log.ScryfallId, log.FromContainer.ContainerId}

			if delta, ok := changes[key]; ok {
				changes[key] = transactionDelta{
					transactionId: delta.transactionId,
					delta:         delta.delta + log.Quantity,
				}
			} else {
				changes[key] = transactionDelta{
					transactionId: log.TransactionId,
					delta:         log.Quantity,
				}
			}
		}
		if log.ToContainer != nil {
			key := changeKey{log.ScryfallId, log.ToContainer.ContainerId}

			if delta, ok := changes[key]; ok {
				changes[key] = transactionDelta{
					transactionId: delta.transactionId,
					delta:         delta.delta - log.Quantity,
				}
			} else {
				changes[key] = transactionDelta{
					transactionId: log.TransactionId,
					delta:         -log.Quantity,
				}
			}
		}
	}

	logsById := make(map[int]*TransactionLogs, len(logs))
	for _, log := range logs {
		logsById[log.TransactionId] = &log
	}

	updatedLogs := make([]TransactionLogs, len(changes))
	for _, delta := range changes {
		log, ok := logsById[delta.transactionId]
		if !ok {
			return logs
		}
		if log == nil {
			return logs
		}

		var newLog TransactionLogs
		if delta.delta > 0 {
			newLog = TransactionLogs{}
		} else {
			newLog = TransactionLogs{
				TransactionId: delta.transactionId,
			}
		}

	}

	additions := map[uuid.UUID][]*TransactionLogs{}
	removals := map[uuid.UUID][]*TransactionLogs{}

	for _, log := range logs {
		if log.FromContainer != nil {
			cardId := log.ScryfallId
			removals[cardId] = append(removals[cardId], &log)
		}
		if log.ToContainer != nil {
			cardId := log.ScryfallId
			additions[cardId] = append(additions[cardId], &log)
		}
	}

	logUpdates := map[int]int{}
	for cardId, adds := range additions {
		removes := removals[cardId]
		i, j := 0, 0
		for i < len(adds) && j < len(removes) {
			add := adds[i]
			if add == nil {
				i += 1
				continue
			}

			remove := removes[j]
			if remove == nil {
				j += 1
				continue
			}

			if add.Quantity < remove.Quantity {
				logUpdates[add.TransactionId] = logUpdates[add.TransactionId]
			} else if add.Quantity > remove.Quantity {
			} else {
			}
			// TODO add some log updates based on joins
		}
	}

	if len(logUpdates) == 0 {
		return logs
	}

	updatedLogs := make([]TransactionLogs, len(logs))
	for _, log := range logs {
		updatedAmount := log.Quantity + logUpdates[log.TransactionId]
		if updatedAmount < 0 {
			continue
		}
		updatedLogs = append(updatedLogs, TransactionLogs{
			log.TransactionId,
			log.GroupId,
			log.FromContainer,
			log.ToContainer,
			log.ScryfallId,
			updatedAmount,
		})
	}

	return updatedLogs
}
