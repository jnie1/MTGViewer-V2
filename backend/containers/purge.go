package containers

import (
	"cmp"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func TranslatePrune(candidates []CardDeposit, cheapCards []cards.CardPricePreview, targetCopies int, minPrice float64) []ContainerChanges {
	if targetCopies < 0 {
		return nil
	}

	cheapPrices := make(map[uuid.UUID]float64, len(cheapCards))
	for _, price := range cheapCards {
		cheapPrices[price.ScryfallId] = price.Price
	}

	cheapOracles := map[uuid.UUID][]CardDeposit{}

	for _, deposit := range candidates {
		cardPrice, ok := cheapPrices[deposit.ScryfallId]
		if !ok {
			continue
		}

		oracleId := deposit.OracleId

		if cardPrice < minPrice {
			cheapOracles[oracleId] = append(cheapOracles[oracleId], deposit)
		}
	}

	removesByContainer := map[int][]CardRequest{}

	for _, deposits := range cheapOracles {
		remainingCopies := targetCopies

		slices.SortFunc(deposits, func(a, b CardDeposit) int {
			priceA, ok := cheapPrices[a.ScryfallId]
			if !ok {
				priceA = minPrice
			}
			priceB, ok := cheapPrices[b.ScryfallId]
			if !ok {
				priceB = minPrice
			}

			priceCmp := cmp.Compare(priceA, priceB)
			if priceCmp != 0 {
				return -priceCmp
			}

			// negative to sort in desc order
			return -cmp.Compare(a.Amount, b.Amount)
		})

		for _, deposit := range deposits {
			cardPrice, ok := cheapPrices[deposit.ScryfallId]

			// assume missing cards are above price, so keep them all
			if cardPrice >= minPrice || !ok {
				replacing := min(deposit.Amount, remainingCopies)
				remainingCopies -= replacing
				continue
			}

			keeping := min(deposit.Amount, remainingCopies)
			removing := deposit.Amount - keeping

			if removing > 0 {
				containerId := deposit.ContainerId
				request := CardRequest{deposit.ScryfallId, deposit.OracleId, -removing}
				removesByContainer[containerId] = append(removesByContainer[containerId], request)
			}

			remainingCopies -= keeping
		}
	}

	removals := make([]ContainerChanges, len(removesByContainer))
	i := 0
	for containerId, requests := range removesByContainer {
		removals[i] = ContainerChanges{containerId, MergeCardRequests(requests)}
		i += 1
	}

	return removals
}

func PreviewPrune(boxes []ContainerPreview, changes []ContainerChanges, fullCards []cards.Card, prices []cards.CardPricePreview) PrunePreviews {
	containerIdRequests := make(map[int][]CardRequest, len(changes))

	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}
	bxIdToName := make(map[int]string, len(boxes))
	for _, box := range boxes {
		bxIdToName[box.ContainerId] = box.Name
	}
	pricesByCard := make(map[uuid.UUID]float64, len(prices))
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	for _, container := range changes {
		var allRequests []CardRequest
		for _, req := range container.Requests {
			allRequests = append(allRequests, req)
		}
		containerIdRequests[container.ContainerId] = allRequests
	}
	grandTotal := 0
	var globalPrnPreviews PrunePreviews
	for containerId, requests := range containerIdRequests {
		var previewCards []cards.CardPriceAmount
		total := 0
		for _, req := range MergeCardRequests(requests) {
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
		grandTotal += total
		globalPrnPreviews.ContainersPrunePreviews = append(globalPrnPreviews.ContainersPrunePreviews, ContainersPrunePreview{
			ContainerId:   containerId,
			ContainerName: bxIdToName[containerId],
			Total:         total,
			Cards:         previewCards,
		})
	}
	slices.SortFunc(globalPrnPreviews.ContainersPrunePreviews, func(a, b ContainersPrunePreview) int {
		return cmp.Compare(a.ContainerId, b.ContainerId)
	})
	globalPrnPreviews.Total = grandTotal
	return globalPrnPreviews
}

func PruneRequests(candidates []CardRequest, fullCards []cards.Card, cheapCards []cards.CardPricePreview, targetCopies int, minPrice float64) []CardRequest {
	if targetCopies < 0 {
		return nil
	}

	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}

	cheapPrices := make(map[uuid.UUID]float64, len(cheapCards))
	for _, price := range cheapCards {
		cheapPrices[price.ScryfallId] = price.Price
	}

	cheapOracles := map[uuid.UUID][]CardRequest{}

	for _, req := range candidates {
		if req.Delta <= 0 {
			continue
		}

		card, ok := cardsById[req.ScryfallId]
		if !ok {
			continue
		}

		cardPrice, ok := cheapPrices[req.ScryfallId]
		if !ok {
			continue
		}

		oracleId := req.OracleId

		if strings.Contains(card.Type, "Basic") || cardPrice < minPrice {
			cheapOracles[oracleId] = append(cheapOracles[oracleId], req)
		}
	}

	var keptCards []CardRequest

	for _, reqs := range cheapOracles {
		remainingCopies := targetCopies

		slices.SortFunc(reqs, func(a, b CardRequest) int {
			priceA, ok := cheapPrices[a.ScryfallId]
			if !ok {
				priceA = minPrice
			}
			priceB, ok := cheapPrices[b.ScryfallId]
			if !ok {
				priceB = minPrice
			}

			priceCmp := cmp.Compare(priceA, priceB)
			if priceCmp != 0 {
				return -priceCmp
			}

			// negative to sort in desc order
			return -cmp.Compare(a.Delta, b.Delta)
		})

		for _, req := range reqs {
			card, ok := cardsById[req.ScryfallId]
			if !ok {
				continue
			}

			if strings.Contains(card.Type, "Basic") {
				continue
			}

			cardPrice, ok := cheapPrices[req.ScryfallId]
			if !ok {
				continue
			}

			if cardPrice >= minPrice {
				replacing := min(req.Delta, remainingCopies)
				remainingCopies -= replacing
				continue
			}

			keeping := min(req.Delta, remainingCopies)

			if keeping > 0 {
				req := CardRequest{req.ScryfallId, req.OracleId, keeping}
				keptCards = append(keptCards, req)
			}

			remainingCopies -= keeping
		}
	}

	for _, req := range candidates {
		card, ok := cardsById[req.ScryfallId]
		if !ok {
			continue
		}

		if strings.Contains(card.Type, "Basic") {
			continue
		}

		// if not cheap, then assume it's expensive, so keep
		if _, ok := cheapPrices[req.ScryfallId]; !ok {
			keptCards = append(keptCards, req)
		}
	}

	return MergeCardRequests(keptCards)

}
