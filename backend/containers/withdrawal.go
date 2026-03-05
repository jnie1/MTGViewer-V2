package containers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

var ErrNegativeWithdrawal = errors.New("negative withdrawal amount specified")
var ErrInsufficientDeposits = errors.New("unsufficient cards in containers to fullfill withdrawal")

func ResolveExtraIdentifiers(withdrawals ContainerWithdrawals) error {
	identifierOptions := map[int]cards.CardIdentifier{}
	extraIds := []cards.CardIdentifier{}

	for _, targets := range withdrawals {
		for _, target := range targets {
			if target.Card == nil {
				return ErrUnknownCardIdentifier
			}

			if _, ok := target.Card.(ScryfallIdentifier); ok {
				continue
			}

			copy := target.Card.Copy()

			var key int
			switch copy.(type) {
			case cards.NameSet:
				key = 3
			case cards.SetCollectorNumber:
				key = 2
			case cards.MultiverseIdentifier:
			case cards.ScryfallIdentifier:
			default:
				key = 1
			}

			// last one wins, intentional
			identifierOptions[key] = copy
			extraIds = append(extraIds, copy)
		}
	}

	if len(extraIds) == 0 {
		return nil
	}

	results, err := cards.FetchCollection(extraIds)
	if err != nil {
		return err
	}

	// hacky way to remap cards back to source id, just get all unique id types per card
	scryfallIdMappings := make(map[cards.CardIdentifier]uuid.UUID, len(results)*len(identifierOptions))

	for _, card := range results {
		for _, converter := range identifierOptions {
			if id, err := converter.Convert(card); err == nil {
				scryfallIdMappings[id] = card.ScryfallId
			}
		}
	}

	for _, targets := range withdrawals {
		for i := range targets {
			target := targets[i]
			if target.Card == nil {
				return ErrUnknownCardIdentifier
			}
			copy := target.Card.Copy()
			if scryfallId, ok := scryfallIdMappings[copy]; ok {
				targets[i] = CardIdentifierAmount{
					Card:   ScryfallIdentifier{Id: scryfallId},
					Amount: target.Amount,
				}
			}
		}
	}

	return nil
}

type depositKey struct {
	ContainerId int
	ScryfallId  uuid.UUID
}

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

			scryfallId, ok := withdrawal.Card.(ScryfallIdentifier)
			if !ok {
				return nil, ErrUnknownCardIdentifier
			}

			key := depositKey{containerId, scryfallId.Id}
			if amountsByContainers[key]-withdrawal.Amount < 0 {
				return nil, ErrInsufficientDeposits
			}

			requests = append(requests, CardRequest{scryfallId.Id, -withdrawal.Amount})
		}

		changes = append(changes, ContainerChanges{containerId, requests})
	}

	return changes, nil
}

func FindScryfallIds(withdrawals ContainerWithdrawals) uuid.UUIDs {
	uniqIds := map[ScryfallIdentifier]bool{}
	for _, targets := range withdrawals {
		for _, target := range targets {
			if scryfallId, ok := target.Card.(ScryfallIdentifier); ok {
				uniqIds[scryfallId] = true
			}
		}
	}
	identifiers := make(uuid.UUIDs, len(uniqIds))
	i := 0
	for scryfallId := range uniqIds {
		identifiers[i] = scryfallId.Id
		i += 1
	}
	return identifiers
}
