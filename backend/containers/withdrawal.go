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
	var unknownOracle cards.ScryfallOracleObj

	scryfallIds := map[uuid.UUID]cards.ScryfallOracleObj{}
	multiverseIds := map[int]cards.ScryfallOracleObj{}

	setCollectors := map[cards.SetCollectorNumber]cards.ScryfallOracleObj{}
	nameSets := map[cards.NameSet]cards.ScryfallOracleObj{}

	for _, targets := range withdrawals {
		for i := range targets {
			switch t := targets[i].Card.(type) {
			case cards.ScryfallIdObj:
				scryfallIds[t.ScryfallId] = unknownOracle

			case cards.MultiverseIdObj:
				multiverseIds[t.MultiverseId] = unknownOracle

			case cards.SetCollectorNumber:
				setCollectors[t] = unknownOracle

			case cards.NameSet:
				nameSets[t] = unknownOracle

			case cards.ScryfallOracleObj:
			default:
				return cards.ErrUnknownCardIdentifier
			}
		}
	}

	if len(scryfallIds) > 0 {
		keys := slices.Collect(maps.Keys(scryfallIds))
		fullCards, err := cards.FetchCollection(ctx, keys...)
		if err != nil {
			return err
		}
		for _, card := range fullCards {
			scryfallIds[card.ScryfallId] = cards.ScryfallOracleObj{
				ScryfallId: card.ScryfallId,
				OracleId:   card.OracleId,
			}
		}
	}

	if len(multiverseIds) > 0 {
		keys := slices.Collect(maps.Keys(multiverseIds))
		cardIds, err := cards.FetchIdsByMultiverseId(ctx, keys...)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			multiverseIds[card.MultiverseId] = cards.ScryfallOracleObj{
				ScryfallId: card.ScryfallId,
				OracleId:   card.OracleId,
			}
		}
	}

	if len(setCollectors) > 0 {
		keys := slices.Collect(maps.Keys(setCollectors))
		cardIds, err := cards.FetchIdsBySetCollector(ctx, keys...)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			setCollectors[card.SetCollectorNumber()] = cards.ScryfallOracleObj{
				ScryfallId: card.ScryfallId,
				OracleId:   card.OracleId,
			}
		}

	}

	if len(nameSets) > 0 {
		keys := slices.Collect(maps.Keys(nameSets))
		cardIds, err := cards.FetchIdsByNameSet(ctx, keys...)
		if err != nil {
			return err
		}
		for _, card := range cardIds {
			nameSets[card.NameSet()] = cards.ScryfallOracleObj{
				ScryfallId: card.ScryfallId,
				OracleId:   card.OracleId,
			}
		}
	}

	for _, targets := range withdrawals {
		for i, target := range targets {
			switch t := target.Card.(type) {
			case cards.ScryfallIdObj:
				if obj := scryfallIds[t.ScryfallId]; obj != unknownOracle {
					targets[i] = CardIdentifierAmount{obj, target.Amount}
				}
			case cards.MultiverseIdObj:
				if obj := multiverseIds[t.MultiverseId]; obj != unknownOracle {
					targets[i] = CardIdentifierAmount{obj, target.Amount}
				}
			case cards.SetCollectorNumber:
				if obj := setCollectors[t]; obj != unknownOracle {
					targets[i] = CardIdentifierAmount{obj, target.Amount}
				}
			case cards.NameSet:
				if obj := nameSets[t]; obj != unknownOracle {
					targets[i] = CardIdentifierAmount{obj, target.Amount}
				}
			case cards.ScryfallOracleObj:
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

			obj, ok := withdrawal.Card.(cards.ScryfallOracleObj)
			if !ok {
				return nil, ErrExpectedScryfallId
			}

			key := ContainerCard{containerId, obj.ScryfallId}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, ErrInsufficientDeposits
			}

			requests = append(requests, CardRequest{obj.ScryfallId, obj.OracleId, -withdrawal.Amount})
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
			case cards.ScryfallOracleObj:
				uniqIds[t.ScryfallId] = nil
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
