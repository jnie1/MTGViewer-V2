package containers

import (
	"errors"

	"github.com/google/uuid"
)

// have finer grain of control for removing cards
// so can assume request will specify the container changes
// all we have to do is just validate them
// also handle concurrency cases, which maybe we just fail if that happens? or yagni

type depositKey struct {
	ContainerId int
	ScryfallId  uuid.UUID
}

var NegativeWithdrawalError = errors.New("negative withdrawal amount specified")
var InsufficientDepositsError = errors.New("unsufficient cards in containers to fullfill withdrawal")

func ValidateCardWithdrawals(withdrawals ContainerWithdrawals, deposits []CardDeposit) ([]ContainerChanges, error) {
	changes := []ContainerChanges{}
	amountsByContainers := map[depositKey]int{}

	for _, deposit := range deposits {
		key := depositKey{deposit.ContainerId, deposit.ScryfallId}
		amountsByContainers[key] = deposit.Amount
	}

	for containerId, withdrawalTargets := range withdrawals {
		requests := []CardRequest{}

		for _, withdrawal := range withdrawalTargets {
			if withdrawal.Amount < 0 {
				return nil, NegativeWithdrawalError
			}
			key := depositKey{containerId, withdrawal.ScryfallId}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, InsufficientDepositsError
			}

			requests = append(requests, CardRequest{withdrawal.ScryfallId, -withdrawal.Amount})
		}

		changes = append(changes, ContainerChanges{containerId, requests})
	}

	return changes, nil
}
