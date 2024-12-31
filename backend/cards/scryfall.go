package cards

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

var scryfallUrl = "https://api.scryfall.com"

func FetchRandomCard() (Card, error) {
	var card Card

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

	var result scryfallCard
	if err := json.Unmarshal(content, &result); err != nil {
		return card, err
	}

	card = toCard(result)
	return card, nil
}

func FetchCollection() ([]Card, error) {
	var cards []Card

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

	var result collectionResult
	if err := json.Unmarshal(content, &result); err != nil {
		return cards, err
	}

	cards = toCards(result.Cards)
	return cards, nil
}
