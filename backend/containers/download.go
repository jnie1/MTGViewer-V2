package containers

import (
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func WriteCardRequests(writer io.Writer, reqs []CardRequest, fullCards []cards.Card) error {
	cardsById := make(map[uuid.UUID]cards.Card, len(fullCards))
	for _, card := range fullCards {
		cardsById[card.ScryfallId] = card
	}

	for _, req := range reqs {
		if card, ok := cardsById[req.ScryfallId]; ok {
			_, err := fmt.Fprintf(writer, "%d %s (%s) %s", req.Delta, card.Name, card.SetCode, card.CollectorNumber)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
