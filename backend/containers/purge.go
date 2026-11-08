package containers

import (
	"cmp"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func TranslatePrune(options []CardDeposit, prices []cards.CardPricePreview, targetCopies int, minPrice float64) []ContainerChanges {
	if targetCopies < 0 {
		return nil
	}

	pricesByCard := make(map[uuid.UUID]float64, len(prices))
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	cheapOracles := map[uuid.UUID]bool{}

	for _, deposit := range options {
		cardPrice, ok := pricesByCard[deposit.ScryfallId]
		if !ok {
			continue
		}

		if cardPrice < minPrice {
			cheapOracles[deposit.OracleId] = true
		}
	}

	depositsByOracle := map[uuid.UUID][]CardDeposit{}

	for _, deposit := range options {
		oracleId := deposit.OracleId
		if cheapOracles[oracleId] {
			depositsByOracle[oracleId] = append(depositsByOracle[oracleId], deposit)
		}
	}

	changesByContainer := map[int][]CardRequest{}

	for _, deposits := range depositsByOracle {
		remainingCopies := targetCopies

		slices.SortFunc(deposits, func(a, b CardDeposit) int {
			priceA, ok := pricesByCard[a.ScryfallId]
			if !ok {
				priceA = minPrice
			}
			priceB, ok := pricesByCard[b.ScryfallId]
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
			cardPrice, ok := pricesByCard[deposit.ScryfallId]

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

func PruneRequests(options []CardRequest, fullCards []cards.Card, prices []cards.CardPricePreview, targetCopies int, minPrice float64) []CardRequest {
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

	cheapOracles := map[uuid.UUID]bool{}

	for _, req := range options {
		if req.Delta <= 0 {
			continue
		}

		card, ok := cardsById[req.ScryfallId]
		if !ok {
			continue
		}

		cardPrice, ok := pricesByCard[req.ScryfallId]
		if !ok {
			continue
		}

		if strings.Contains(card.Type, "Basic") || cardPrice < minPrice {
			cheapOracles[req.OracleId] = true
		}
	}

	reqsByOracle := map[uuid.UUID][]CardRequest{}

	for _, req := range options {
		oracleId := req.OracleId
		if cheapOracles[oracleId] {
			reqsByOracle[oracleId] = append(reqsByOracle[oracleId], req)
		}
	}

	var result []CardRequest

	for _, reqs := range reqsByOracle {
		remainingCopies := targetCopies

		slices.SortFunc(reqs, func(a, b CardRequest) int {
			priceA, ok := pricesByCard[a.ScryfallId]
			if !ok {
				priceA = minPrice
			}
			priceB, ok := pricesByCard[b.ScryfallId]
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
			cardPrice, ok := pricesByCard[req.ScryfallId]

			// assume missing cards are above price, so keep them all
			if cardPrice >= minPrice || !ok {
				replacing := min(req.Delta, remainingCopies)
				result = append(result, req)
				remainingCopies -= replacing
				continue
			}

			keeping := min(req.Delta, remainingCopies)

			if keeping > 0 {
				req := CardRequest{req.ScryfallId, req.OracleId, keeping}
				result = append(result, req)
			}

			remainingCopies -= keeping
		}
	}

	return MergeCardRequests(result)

}
