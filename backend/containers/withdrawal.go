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
	multiverseIds := map[int]uuid.UUID{}
	setCollectors := map[cards.SetCollectorNumber]uuid.UUID{}
	nameSets := map[cards.NameSet]uuid.UUID{}

	for _, targets := range withdrawals {
		for i := range targets {
			switch t := targets[i].Card.(type) {
			case cards.MultiverseIdObj:
				multiverseIds[t.MultiverseId] = uuid.Nil

			case cards.SetCollectorNumber:
				setCollectors[t] = uuid.Nil

			case cards.NameSet:
				nameSets[t] = uuid.Nil

			case cards.ScryfallIdObj:
			default:
				return cards.ErrUnknownCardIdentifier
			}
		}
	}

	if len(multiverseIds) > 0 {
		keys := slices.Collect(maps.Keys(multiverseIds))
		cardIds, err := cards.FetchIdsByMultiverseId(ctx, keys)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			multiverseIds[card.MultiverseId] = card.ScryfallId
		}
	}

	if len(setCollectors) > 0 {
		keys := slices.Collect(maps.Keys(setCollectors))
		cardIds, err := cards.FetchIdsBySetCollector(ctx, keys)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			setCollectors[card.SetCollectorNumber()] = card.ScryfallId
		}

	}

	if len(nameSets) > 0 {
		keys := slices.Collect(maps.Keys(nameSets))
		cardIds, err := cards.FetchIdsByNameSet(ctx, keys)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			nameSets[card.NameSet()] = card.ScryfallId
		}
	}

	for _, targets := range withdrawals {
		for i, target := range targets {
			switch t := target.Card.(type) {
			case cards.MultiverseIdObj:
				if scryfallId := multiverseIds[t.MultiverseId]; scryfallId != uuid.Nil {
					targets[i] = CardIdentifierAmount{scryfallId, target.Amount}
				}
			case cards.SetCollectorNumber:
				if scryfallId := setCollectors[t]; scryfallId != uuid.Nil {
					targets[i] = CardIdentifierAmount{scryfallId, target.Amount}
				}
			case cards.NameSet:
				if scryfallId := nameSets[t]; scryfallId != uuid.Nil {
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
