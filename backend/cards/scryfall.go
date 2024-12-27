package cards

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ScryfallImages struct {
	Small  string `json:"small"`
	Normal string `json:"normal"`
	Large  string `json:"large"`
}

type ScryfallCard struct {
	ScryfallId    string         `json:"id"`
	MultiverseIds []int          `json:"multiverse_ids"`
	ManaCost      string         `json:"mana_cost"`
	Name          string         `json:"name"`
	Power         string         `json:"power"`
	Toughness     string         `json:"toughness"`
	Images        ScryfallImages `json:"image_uris"`
	Type          string         `json:"type_line"`
}

var scryfallUrl = "https://api.scryfall.com"

func FetchRandomCard() (ScryfallCard, error) {
	card := ScryfallCard{}
	randomUrl, err := url.JoinPath(scryfallUrl, "/cards/random")
	if err != nil {
		return card, err
	}

	req, err := http.NewRequest("GET", randomUrl, nil)
	if err != nil {
		return card, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return card, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return card, err
	}

	err = json.Unmarshal(content, &card)

	return card, err
}

type setCollectorNumber struct {
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
}

type collectionQuery struct {
	Identifiers []setCollectorNumber `json:"identifiers"`
}

type collectionResult struct {
	Cards []ScryfallCard `json:"data"`
}

func FetchCollection() ([]ScryfallCard, error) {
	cards := []ScryfallCard{}

	collectionUrl, err := url.JoinPath(scryfallUrl, "/cards/collection")
	if err != nil {
		return cards, err
	}

	query := collectionQuery{
		Identifiers: []setCollectorNumber{
			{Set: "FDN", CollectorNumber: "100"},
			{Set: "FDN", CollectorNumber: "105"},
		},
	}

	payload, err := json.Marshal(query)
	if err != nil {
		return cards, err
	}

	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", collectionUrl, body)
	if err != nil {
		return cards, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return cards, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return cards, err
	}

	result := collectionResult{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return cards, err
	}

	cards = result.Cards
	return cards, nil
}
