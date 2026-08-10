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
var ErrExpectedScryfallId = errors.New("expected scryfall id uuid")
var ErrInsufficientDeposits = errors.New("unsufficient cards in containers to fullfill request")

func ResolveIdentifiers(ctx context.Context, withdrawals ContainerWithdrawals) error {
	extraQuery, err := FindIdQuery(withdrawals)
	if err != nil {
		return err
	}

	multiverseIds := map[int]uuid.UUID{}
	setNumbers := map[cards.SetCollectorNumber]uuid.UUID{}
	nameSets := map[cards.NameSet]uuid.UUID{}

	if !extraQuery.IsEmpty() {
		extraIds, err := cards.FetchIdentifiers(ctx, extraQuery)
		if err != nil {
			return err
		}
		for _, card := range extraIds {
			multiverseIds[card.MultiverseId] = card.ScryfallId
			setNumbers[card.SetCollectorNumber()] = card.ScryfallId
			nameSets[card.NameSet()] = card.ScryfallId
		}
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

func ValidateCardWithdrawals(withdrawals ContainerWithdrawals, deposits []CardDepositPreview) ([]ContainerChanges, error) {
	changes := []ContainerChanges{}
	amountsByContainers := map[ContainerCard]int{}

	for _, deposit := range deposits {
		key := ContainerCard{deposit.ContainerId, deposit.ScryfallId}
		amountsByContainers[key] = deposit.Amount
	}

	for containerId, targets := range withdrawals {
		requests := []CardRequest{}

		for _, withdrawal := range targets {
			if withdrawal.Amount < 0 {
				return nil, ErrNegativeWithdrawal
			}

			id, ok := withdrawal.Card.(uuid.UUID)
			if !ok {
				return nil, ErrExpectedScryfallId
			}

			key := ContainerCard{containerId, id}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, ErrInsufficientDeposits
			}

			requests = append(requests, CardRequest{id, -withdrawal.Amount})
		}

		changes = append(changes, ContainerChanges{containerId, requests})
	}

	return changes, nil
}

func FindIdQuery(withdrawals ContainerWithdrawals) (cards.CardIdQuery, error) {
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
