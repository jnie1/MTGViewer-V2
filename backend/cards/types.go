package cards

type CardImageUrls struct {
	Preview string `json:"preview"`
	Normal  string `json:"normal"`
	Full    string `json:"full"`
}

type Card struct {
	Name      string        `json:"name"`
	ManaCost  string        `json:"manaCost"`
	Type      string        `json:"type"`
	Power     string        `json:"power"`
	Toughness string        `json:"toughness"`
	Images    CardImageUrls `json:"imageUrls"`
}

type scryfallImages struct {
	Small  string `json:"small"`
	Normal string `json:"normal"`
	Large  string `json:"large"`
}

type scryfallCard struct {
	ScryfallId string         `json:"id"`
	ManaCost   string         `json:"mana_cost"`
	Name       string         `json:"name"`
	Power      string         `json:"power"`
	Toughness  string         `json:"toughness"`
	Images     scryfallImages `json:"image_uris"`
	Type       string         `json:"type_line"`
}

type setCollectorNumber struct {
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
}

type collectionQuery struct {
	Identifiers []setCollectorNumber `json:"identifiers"`
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
		Power:     card.Power,
		Toughness: card.Toughness,
		Images: CardImageUrls{
			Preview: images.Small,
			Normal:  images.Normal,
			Full:    images.Large,
		},
	}
}

func toCards(cards []scryfallCard) []Card {
	result := make([]Card, len(cards))
	for i, card := range cards {
		result[i] = toCard(card)
	}
	return result
}
