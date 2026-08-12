package cards

import (
	"strings"

	"github.com/google/uuid"
)

func FindCheapCards(cards []Card, amounts []CardAmountPreview, prices []CardPrice, price float64) []CardAmount {
	belowPrice := []CardAmount{}

	amountsByCard := map[uuid.UUID]int{}
	for _, amt := range amounts {
		amountsByCard[amt.ScryfallId] = amt.Amount
	}

	pricesByCard := map[uuid.UUID]float64{}
	for _, price := range prices {
		pricesByCard[price.ScryfallId] = price.Price
	}

	for _, card := range cards {
		if strings.Contains(card.Type, "Land") {
			continue
		}
		cardPrice := pricesByCard[card.ScryfallId]
		if cardPrice <= price {
			belowPrice = append(belowPrice, CardAmount{card, amountsByCard[card.ScryfallId]})
		}
	}

	return belowPrice
}
