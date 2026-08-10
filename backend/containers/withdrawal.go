package containers

import (
	"context"
	"errors"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

var ErrNegativeWithdrawal = errors.New("negative withdrawal amount specified")
var ErrInsufficientDeposits = errors.New("unsufficient cards in containers to fullfill request")

func ResolveExtraIdentifiers(ctx context.Context, withdrawals ContainerWithdrawals) error {
	extraQuery, err := FindIdentifiers(withdrawals)
	if err != nil {
		return err
	}

	if !extraQuery.IsEmpty() {
		return nil
	}

	results, err := cards.FetchIdentifiers(ctx, extraQuery)
	if err != nil {
		return err
	}

	multiverseIds := map[int]uuid.UUID{}
	setNumbers := map[cards.SetCollectorNumber]uuid.UUID{}
	nameSets := map[cards.NameSet]uuid.UUID{}

	for _, card := range results {
		multiverseIds[card.MultiverseId] = card.ScryfallId
		setNumbers[card.SetCollectorNumber()] = card.ScryfallId
		nameSets[card.NameSet()] = card.ScryfallId
	}

	for _, targets := range withdrawals {
		for i, target := range targets {
			switch t := target.Card.(type) {
			case cards.MultiverseIdObj:
				if scryfallId, ok := multiverseIds[t.MultiverseId]; ok {
					targets[i] = CardIdentifierAmount{scryfallId, target.Amount}
				}
			case cards.SetCollectorNumber:
				if scryfallId, ok := setNumbers[t]; ok {
					targets[i] = CardIdentifierAmount{scryfallId, target.Amount}
				}
			case cards.NameSet:
				if scryfallId, ok := nameSets[t]; ok {
					targets[i] = CardIdentifierAmount{scryfallId, target.Amount}
				}
			case cards.ScryfallIdObj:
				targets[i] = CardIdentifierAmount{t.ScryfallId, target.Amount}
			default:
				return cards.ErrUnknownCardIdentifier
			}

		}
	}

	return nil
}

type depositKey struct {
	ContainerId int
	ScryfallId  uuid.UUID
}

func ValidateCardWithdrawals(withdrawals ContainerWithdrawals, deposits []CardDepositPreview) ([]ContainerChanges, error) {
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

			obj, ok := withdrawal.Card.(cards.ScryfallIdObj)
			if !ok {
				return nil, cards.ErrUnknownCardIdentifier
			}

			key := depositKey{containerId, obj.ScryfallId}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, ErrInsufficientDeposits
			}

			requests = append(requests, CardRequest{obj.ScryfallId, -withdrawal.Amount})
		}

		changes = append(changes, ContainerChanges{containerId, requests})
	}

	return changes, nil
}

func FindIdentifiers(withdrawals ContainerWithdrawals) (cards.CardIdQuery, error) {
	multiverseIds := map[int]any{}
	nameSets := map[cards.NameSet]any{}
	collectorNumbers := map[cards.SetCollectorNumber]any{}

	for _, targets := range withdrawals {
		for _, target := range targets {
			switch t := target.Card.(type) {
			case cards.MultiverseIdObj:
				multiverseIds[t.MultiverseId] = nil

			case cards.NameSet:
				nameSets[t] = nil

			case cards.SetCollectorNumber:
				collectorNumbers[t] = nil

			case cards.ScryfallIdObj:
			default:
				return cards.CardIdQuery{}, cards.ErrUnknownCardIdentifier
			}
		}
	}

	ids := cards.CardIdQuery{
		MultiverseIds: slices.Collect(maps.Keys(multiverseIds)),
		SetNumbers:    slices.Collect(maps.Keys(collectorNumbers)),
		NameSets:      slices.Collect(maps.Keys(nameSets)),
	}

	return ids, nil
}

func FindScryfallIds(withdrawals ContainerWithdrawals) uuid.UUIDs {
	uniqIds := map[uuid.UUID]any{}
	for _, targets := range withdrawals {
		for _, target := range targets {
			switch t := target.Card.(type) {
			case cards.ScryfallIdObj:
				uniqIds[t.ScryfallId] = nil
			case uuid.UUID:
				uniqIds[t] = nil
			}
		}
	}

	identifiers := make(uuid.UUIDs, len(uniqIds))
	i := 0
	for id := range uniqIds {
		identifiers[i] = id
		i += 1
	}

	return identifiers
}
