package cards

import (
	"strconv"
	"strings"
)

func KeepBelowPrice(allCards []Card, price float64) ([]Card, error) {
	belowPrice := []Card{}
	for _, card := range allCards {
		cardPrice, err := strconv.ParseFloat(card.Prices["usd"], 64)
		if err != nil {
			return nil, err
		}
		if cardPrice <= price && !strings.Contains(card.Type, "Land") {
			belowPrice = append(belowPrice, card)
		}
	}
	return belowPrice, nil
}
