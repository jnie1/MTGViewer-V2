package cards

type CardImageUrls struct {
	Preview string `json:"preview,omitempty"`
	Normal  string `json:"normal,omitempty"`
	Full    string `json:"full,omitempty"`
}

type Card struct {
	Name      string        `json:"name"`
	ManaCost  string        `json:"manaCost,omitempty"`
	Type      string        `json:"type"`
	Rarity    string        `json:"rarity"`
	Power     string        `json:"power,omitempty"`
	Toughness string        `json:"toughness,omitempty"`
	Images    CardImageUrls `json:"imageUrls"`
}

type scryfallImages struct {
	Small  string `json:"small,omitempty"`
	Normal string `json:"normal,omitempty"`
	Large  string `json:"large,omitempty"`
}

type scryfallCard struct {
	ScryfallId string         `json:"id"`
	ManaCost   string         `json:"mana_cost,omitempty"`
	Name       string         `json:"name"`
	Power      string         `json:"power,omitempty"`
	Toughness  string         `json:"toughness,omitempty"`
	Images     scryfallImages `json:"image_uris"`
	Type       string         `json:"type_line"`
	Rarity     string         `json:"rarity"`
}

type scryfallIdentifier struct {
	Id string `json:"id"`
}

type collectionQuery struct {
	Identifiers []scryfallIdentifier `json:"identifiers"`
}

type collectionResult struct {
	Cards []scryfallCard `json:"data"`
}

func toCard(card scryfallCard) Card {
	images := card.Images
	return Card{
		Name:      card.Name,
		ManaCost:  card.ManaCost,
		Type:      card.Type,
		Rarity:    card.Rarity,
		Power:     card.Power,
		Toughness: card.Toughness,
		Images: CardImageUrls{
			Preview: images.Small,
			Normal:  images.Normal,
			Full:    images.Large,
		},
	}
}

func toScryfallIdentifiers(ids []string) []scryfallIdentifier {
	result := make([]scryfallIdentifier, len(ids))
	for i, id := range ids {
		result[i] = scryfallIdentifier{Id: id}
	}
	return result
}

func toCards(cards []scryfallCard) []Card {
	result := make([]Card, len(cards))
	for i, card := range cards {
		result[i] = toCard(card)
	}
	return result
}
