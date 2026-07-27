package cards

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func FindCheapCards(cards []Card, amounts []CardAmountPreview, price float64) ([]CardAmount, error) {
	belowPrice := []CardAmount{}
	amountsById := map[uuid.UUID]int{}

	for _, amt := range amounts {
		amountsById[amt.ScryfallId] = amt.Amount
	}

	for _, card := range cards {
		cardPrice, err := strconv.ParseFloat(card.Prices["usd"], 64)
		if err != nil {
			return nil, err
		}
		if cardPrice <= price && !strings.Contains(card.Type, "Land") {
			belowPrice = append(belowPrice, CardAmount{card, amountsById[card.ScryfallId]})
		}
	}
	return belowPrice, nil
}
