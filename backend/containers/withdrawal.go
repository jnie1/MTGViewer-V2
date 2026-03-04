package containers

import (
	"errors"

	"github.com/google/uuid"
)

func FindScryfallIds(withdrawals ContainerWithdrawals) uuid.UUIDs {
	uniqIds := map[uuid.UUID]bool{}
	for _, targets := range withdrawals {
		for _, withdrawal := range targets {
			uniqIds[withdrawal.ScryfallId] = true
		}
	}
	scryfallIds := make(uuid.UUIDs, len(uniqIds))
	i := 0
	for id := range uniqIds {
		scryfallIds[i] = id
		i += 1
	}
	return scryfallIds
}

type depositKey struct {
	ContainerId int
	ScryfallId  uuid.UUID
}

var ErrNegativeWithdrawal = errors.New("negative withdrawal amount specified")
var ErrInsufficientDeposits = errors.New("unsufficient cards in containers to fullfill withdrawal")

func ValidateCardWithdrawals(withdrawals ContainerWithdrawals, deposits []CardDeposit) ([]ContainerChanges, error) {
	changes := []ContainerChanges{}
	amountsByContainers := map[depositKey]int{}

	for _, deposit := range deposits {
		key := depositKey{deposit.ContainerId, deposit.ScryfallId}
		amountsByContainers[key] = deposit.Amount
	}

	for containerId, targets := range withdrawals {
		requests := []CardRequest{}

		for _, withdrawal := range targets {
			if withdrawal.Amount < 0 {
				return nil, ErrNegativeWithdrawal
			}
			key := depositKey{containerId, withdrawal.ScryfallId}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, ErrInsufficientDeposits
			}

			requests = append(requests, CardRequest{withdrawal.ScryfallId, -withdrawal.Amount})
		}

		changes = append(changes, ContainerChanges{containerId, requests})
	}

	return changes, nil
}
