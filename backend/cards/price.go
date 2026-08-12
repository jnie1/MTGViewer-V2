package cards

import (
	"strings"

	"github.com/google/uuid"
)

func FindCheapCards(cards []Card, amounts []CardAmountPreview, prices []CardPricePreview, price float64) []CardPriceAmount {
	belowPrice := []CardPriceAmount{}

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
		cardPrice, ok := pricesByCard[card.ScryfallId]
		if !ok {
			continue
		}

		if cardPrice <= price {
			belowPrice = append(belowPrice, CardPriceAmount{
				CardAmount{card, amountsByCard[card.ScryfallId]}, cardPrice})
		}
	}

	return belowPrice
}
