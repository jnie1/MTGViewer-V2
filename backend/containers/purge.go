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

	changesByCard := map[ContainerCard]int{}
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
				cardKey := ContainerCard{deposit.ContainerId, deposit.ScryfallId}
				changesByCard[cardKey] = changesByCard[cardKey] - pruneAmount
			}
			remainingCopies -= keepAmount
		}
	}

	changesByContainer := map[int][]CardRequest{}
	for containerCard, delta := range changesByCard {
		containerId := containerCard.ContainerId
		request := CardRequest{containerCard.ScryfallId, delta}
		changesByContainer[containerId] = append(changesByContainer[containerId], request)
	}

	changes := make([]ContainerChanges, len(changesByContainer))
	i := 0
	for containerId, requests := range changesByContainer {
		changes[i] = ContainerChanges{containerId, requests}
		i += 1
	}

	return changes
}
