package transactions

import "github.com/google/uuid"

func MergeLogs(logs []TransactionLogs) []TransactionLogs {
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
