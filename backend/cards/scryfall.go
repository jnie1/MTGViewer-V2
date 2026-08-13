package cards

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	scryfallApiUrl   string = "https://api.scryfall.com"
	scryfallImageUrl string = "https://cards.scryfall.io"
)

type scryfallImages struct {
	Small  string `json:"small,omitempty"`
	Normal string `json:"normal,omitempty"`
	Large  string `json:"large,omitempty"`
}

type scryfallCardFace struct {
	Name     string         `json:"name"`
	ManaCost string         `json:"mana_cost,omitempty"`
	Type     string         `json:"type_line"`
	Images   scryfallImages `json:"image_uris"`
}

type scryfallCard struct {
	ScryfallId      uuid.UUID          `json:"id"`
	ManaCost        string             `json:"mana_cost,omitempty"`
	Name            string             `json:"name"`
	SetName         string             `json:"set_name"`
	Set             string             `json:"set"`
	CollectorNumber string             `json:"collector_number"`
	MultiverseIds   []int              `json:"multiverse_ids,omitempty"`
	Power           string             `json:"power,omitempty"`
	Toughness       string             `json:"toughness,omitempty"`
	Images          scryfallImages     `json:"image_uris"`
	CardFaces       []scryfallCardFace `json:"card_faces,omitempty"`
	Type            string             `json:"type_line"`
	Rarity          string             `json:"rarity"`
}

type searchResult struct {
	TotalCards int            `json:"total_cards"`
	HasMore    bool           `json:"has_more"`
	Cards      []scryfallCard `json:"data"`
}

func ImageURLs(scryfallId uuid.UUID) (CardImageURLs, error) {
	preview, err := imageURL(scryfallId, "thumb", "front", "webp")
	if err != nil {
		return CardImageURLs{}, err
	}

	normal, err := imageURL(scryfallId, "grid", "front", "webp")
	if err != nil {
		return CardImageURLs{}, err
	}

	full, err := imageURL(scryfallId, "display", "front", "webp")
	if err != nil {
		return CardImageURLs{}, err
	}

	return CardImageURLs{preview, normal, full}, nil
}

func imageURL(scryfallId uuid.UUID, size, face, ext string) (string, error) {
	fileName := fmt.Sprintf("%s.%s", scryfallId, ext)
	return url.JoinPath(scryfallImageUrl, size, face, string(fileName[0]), string(fileName[1]), fileName)
}

func SearchCards(query string, page int) (SearchCardPage, error) {
	query, err := url.QueryUnescape(query)
	if err != nil {
		return SearchCardPage{}, err
	}

	searchPath, err := url.JoinPath(scryfallApiUrl, "/cards/search")
	if err != nil {
		return SearchCardPage{}, err
	}

	searchUrl, err := url.Parse(searchPath)
	if err != nil {
		return SearchCardPage{}, err
	}

	searchParams := url.Values{}
	searchParams.Add("unique", "prints")
	searchParams.Add("page", strconv.Itoa(page))
	searchParams.Add("q", query)

	searchUrl.RawQuery = searchParams.Encode()
	req, err := http.NewRequest("GET", searchUrl.String(), nil)

	if err != nil {
		return SearchCardPage{}, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SearchCardPage{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return SearchCardPage{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return SearchCardPage{}, fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return SearchCardPage{}, fmt.Errorf("unexpected response content %s", contentType)
	}

	decoder := json.NewDecoder(resp.Body)

	var result searchResult
	if err := decoder.Decode(&result); err != nil {
		return SearchCardPage{}, err
	}

	searchPage := SearchCardPage{
		TotalCards: result.TotalCards,
		Cards:      toCards(result.Cards),
		Page:       page,
		HasMore:    result.HasMore,
	}

	return searchPage, nil
}

func toCard(card scryfallCard) Card {
	var multiverseId int
	if len(card.MultiverseIds) == 1 {
		multiverseId = card.MultiverseIds[0]
	}

	images := card.Images
	if len(card.CardFaces) > 0 && card.CardFaces[0].Images.Small != "" {
		images = card.CardFaces[0].Images
	}

	return Card{
		card.ScryfallId,
		card.Name,
		card.ManaCost,
		card.SetName,
		strings.ToUpper(card.Set),
		card.CollectorNumber,
		multiverseId,
		card.Type,
		card.Rarity,
		card.Power,
		card.Toughness,
		CardImageURLs{
			images.Small,
			images.Normal,
			images.Large,
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
