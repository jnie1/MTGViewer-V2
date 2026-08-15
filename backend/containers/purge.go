package containers

import (
	"cmp"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func TranslatePrune(options []CardDepositPreview, fullCards []cards.Card, prices []cards.CardPricePreview, targetCopies int, minPrice float64) []ContainerChanges {
	if targetCopies < 0 {
		return nil
	}

	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}

	pricesByCard := make(map[uuid.UUID]float64, len(prices))
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	belowPrice := map[uuid.UUID]CardDepositPreview{}
	cheapOracles := map[uuid.UUID]any{}

	for _, deposit := range options {
		card, ok := cardsById[deposit.ScryfallId]
		if !ok {
			continue
		}

		cardPrice, ok := pricesByCard[deposit.ScryfallId]
		if !ok {
			continue
		}

		if strings.Contains(card.Type, "Land") {
			continue
		}

		if cardPrice < minPrice {
			belowPrice[deposit.ScryfallId] = deposit
			cheapOracles[deposit.OracleId] = nil
		}
	}

	depositsByOracle := map[uuid.UUID][]CardDepositPreview{}

	for _, deposit := range options {
		oracleId := deposit.OracleId
		if _, ok := belowPrice[deposit.ScryfallId]; ok {
			depositsByOracle[oracleId] = append(depositsByOracle[oracleId], deposit)
		} else if _, ok := cheapOracles[oracleId]; ok {
			depositsByOracle[oracleId] = append(depositsByOracle[oracleId], deposit)
		}
	}

	changesByContainer := map[int][]CardRequest{}

	for _, deposits := range depositsByOracle {
		remainingCopies := targetCopies

		slices.SortFunc(deposits, func(a, b CardDepositPreview) int {
			priceCmp := cmp.Compare(pricesByCard[a.ScryfallId], pricesByCard[b.ScryfallId])
			if priceCmp != 0 {
				return -priceCmp
			}
			// negative to sort in desc order
			return -cmp.Compare(a.Amount, b.Amount)
		})

		for _, deposit := range deposits {
			cardPrice := pricesByCard[deposit.ScryfallId]

			if cardPrice >= minPrice {
				replacing := min(deposit.Amount, remainingCopies)
				remainingCopies -= replacing
				continue
			}

			keeping := min(deposit.Amount, remainingCopies)
			removing := deposit.Amount - keeping

			if removing > 0 {
				containerId := deposit.ContainerId
				request := CardRequest{deposit.ScryfallId, deposit.OracleId, -removing}
				changesByContainer[containerId] = append(changesByContainer[containerId], request)
			}

			remainingCopies -= keeping
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

func PreviewPrune(changes []ContainerChanges, fullCards []cards.Card, prices []cards.CardPricePreview) CardPrunePreview {
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
			continue
		}

		price, ok := pricesByCard[cardId]
		if !ok {
			continue
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

	return CardPrunePreview{total, previewCards}
}
