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

	changesByContainer := map[int][]CardRequest{}

	for _, deposits := range depositsByCard {
		remainingCopies := targetCopies

		slices.SortFunc(deposits, func(a, b CardDepositPreview) int {
			// negative to sort in desc order
			return -cmp.Compare(a.Amount, b.Amount)
		})

		for _, deposit := range deposits {
			keepAmount := min(deposit.Amount, remainingCopies)
			pruneAmount := deposit.Amount - keepAmount

			if pruneAmount > 0 {
				containerId := deposit.ContainerId
				request := CardRequest{deposit.ScryfallId, -pruneAmount}
				changesByContainer[containerId] = append(changesByContainer[containerId], request)
			}

			remainingCopies -= keepAmount
		}
	}

	changes := make([]ContainerChanges, len(changesByContainer))
	i := 0
	for containerId, requests := range changesByContainer {
		changes[i] = ContainerChanges{containerId, MergeCardRequests(requests)}
		i += 1
	}

	return changes
}
