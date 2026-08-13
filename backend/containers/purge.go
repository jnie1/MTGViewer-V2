package containers

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func FindCheapCandidates(options []cards.Card, prices []cards.CardPricePreview, price float64) uuid.UUIDs {
	pricesByCard := make(map[uuid.UUID]float64, len(prices))
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	belowPrice := uuid.UUIDs{}

	for _, card := range options {
		if strings.Contains(card.Type, "Land") {
			continue
		}
		cardPrice, ok := pricesByCard[card.ScryfallId]
		if !ok {
			continue
		}

		if cardPrice <= price {
			belowPrice = append(belowPrice, card.ScryfallId)
		}
	}

	return belowPrice
}

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
				request := CardRequest{deposit.ScryfallId, deposit.OracleId, -pruneAmount}
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

func PreviewPrune(changes []ContainerChanges, fullCards []cards.Card, prices []cards.CardPricePreview) (CardPrunePreview, error) {
	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}

	pricesByCard := make(map[uuid.UUID]float64, len(prices))
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	allRequests := []CardRequest{}
	for _, container := range changes {
		for _, req := range container.Requests {
			allRequests = append(allRequests, req)
		}
	}

	previewCards := []cards.CardPriceAmount{}
	total := 0

	for _, req := range MergeCardRequests(allRequests) {
		cardId := req.ScryfallId
		amount := req.Delta
		if amount < 0 {
			amount = -amount
		}

		card, ok := cardsById[cardId]
		if !ok {
			return CardPrunePreview{}, fmt.Errorf("unknown scryfall id found: %s", cardId)
		}

		price, ok := pricesByCard[cardId]
		if !ok {
			return CardPrunePreview{}, fmt.Errorf("unknown scryfall id found: %s", cardId)
		}

		cardPrice := cards.CardPriceAmount{
			CardAmount: cards.CardAmount{Card: card, Amount: amount},
			Price:      price,
		}

		previewCards = append(previewCards, cardPrice)
		total += amount
	}

	slices.SortFunc(previewCards, func(a, b cards.CardPriceAmount) int {
		nameCompare := strings.Compare(a.Name, b.Name)
		if nameCompare != 0 {
			return nameCompare
		}
		setCompare := strings.Compare(a.Set, b.Set)
		if setCompare != 0 {
			return setCompare
		}
		return cmp.Compare(a.Price, b.Price)
	})

	return CardPrunePreview{total, previewCards}, nil
}
