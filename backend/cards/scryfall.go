package cards

import (
	"context"
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
	Small   string `json:"small,omitempty"`
	Normal  string `json:"normal,omitempty"`
	Large   string `json:"large,omitempty"`
	Thumb   string `json:"thumb,omitempty"`
	Grid    string `json:"grid,omitempty"`
	Display string `json:"display,omitempty"`
}

type scryfallCardFace struct {
	Name     string         `json:"name"`
	ManaCost string         `json:"mana_cost,omitempty"`
	Type     string         `json:"type_line"`
	Images   scryfallImages `json:"image_uris"`
}

type scryfallCard struct {
	ScryfallId      uuid.UUID          `json:"id"`
	OracleId        uuid.UUID          `json:"oracle_id"`
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
	var imageUrls CardImageURLs

	preview, err := imageURL(scryfallId, "thumb", "front", "webp")
	if err != nil {
		return imageUrls, err
	}

	normal, err := imageURL(scryfallId, "grid", "front", "webp")
	if err != nil {
		return imageUrls, err
	}

	full, err := imageURL(scryfallId, "display", "front", "webp")
	if err != nil {
		return imageUrls, err
	}

	imageUrls.Preview = preview
	imageUrls.Normal = normal
	imageUrls.Full = full

	return imageUrls, nil
}

func imageURL(scryfallId uuid.UUID, size, face, ext string) (string, error) {
	fileName := fmt.Sprintf("%s.%s", scryfallId, ext)
	return url.JoinPath(scryfallImageUrl, size, face, string(fileName[0]), string(fileName[1]), fileName)
}

func SearchCards(ctx context.Context, query string, page int) (SearchCardPage, error) {
	var searchPage SearchCardPage

	query, err := url.QueryUnescape(query)
	if err != nil {
		return searchPage, err
	}

	searchPath, err := url.JoinPath(scryfallApiUrl, "/cards/search")
	if err != nil {
		return searchPage, err
	}

	searchUrl, err := url.Parse(searchPath)
	if err != nil {
		return searchPage, err
	}

	searchParams := url.Values{}
	searchParams.Add("page", strconv.Itoa(page))
	searchParams.Add("q", query)

	searchUrl.RawQuery = searchParams.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchUrl.String(), nil)

	if err != nil {
		return searchPage, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return searchPage, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return searchPage, nil
	}

	if resp.StatusCode != http.StatusOK {
		return searchPage, fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return searchPage, fmt.Errorf("unexpected response content %s", contentType)
	}

	decoder := json.NewDecoder(resp.Body)

	var result searchResult
	if err := decoder.Decode(&result); err != nil {
		return searchPage, err
	}

	searchPage.Cards = toCards(result.Cards)
	searchPage.Page = page
	searchPage.HasMore = result.HasMore

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
		card.OracleId,
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
			images.Thumb,
			images.Grid,
			images.Display,
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
