package containers

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
)

func TranslatePrune(matches []CardDepositPreview, targetCopies int) []ContainerChanges {
	if targetCopies < 0 {
		return nil
	}

	depositsByCard := map[uuid.UUID][]CardDepositPreview{}
	for _, deposit := range matches {
		scryfallId := deposit.ScryfallId
		depositsByCard[scryfallId] = append(depositsByCard[scryfallId], deposit)
	}

	changesByCard := map[uuid.UUID][]ContainerDelta{}
	for _, deposits := range depositsByCard {
		slices.SortFunc(deposits, func(a, b CardDepositPreview) int {
			// negative to sort in desc order
			return -cmp.Compare(a.Amount, b.Amount)
		})

		remainingCopies := targetCopies

		for _, deposit := range deposits {
			keepAmount := min(deposit.Amount, remainingCopies)
			pruneAmount := deposit.Amount - keepAmount

			if pruneAmount > 0 {
				scryfallId := deposit.ScryfallId
				delta := ContainerDelta{deposit.ContainerId, -pruneAmount}
				changesByCard[scryfallId] = append(changesByCard[scryfallId], delta)
			}

			remainingCopies -= keepAmount
		}
	}

	changesByContainer := map[int][]CardRequest{}
	for cardId, changes := range changesByCard {
		for _, change := range changes {
			containerId := change.ContainerId
			request := CardRequest{cardId, change.Delta}
			changesByContainer[containerId] = append(changesByContainer[containerId], request)
		}
	}

	changes := make([]ContainerChanges, len(changesByContainer))
	i := 0
	for containerId, requests := range changesByContainer {
		changes[i] = ContainerChanges{containerId, requests}
		i += 1
	}

	return changes
}
