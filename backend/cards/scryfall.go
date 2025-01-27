package cards

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return card, fmt.Errorf("unexpected response content %s", contentType)
	}

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

type collectionBatchResult struct {
	Cards []Card
	Error error
}

func FetchCollection(scryfallIds []string) ([]Card, error) {
	batchSizeLimit := 75

	results := make(chan collectionBatchResult)
	workerCount := 0

	for batch := range slices.Chunk(scryfallIds, batchSizeLimit) {
		workerCount++
		go func() {
			cards, err := fetchCollectionBatch(batch)
			results <- collectionBatchResult{cards, err}
		}()
	}

	var cards [][]Card
	var errs []error

	for i := 0; i < workerCount; i++ {
		result := <-results
		if result.Error != nil {
			errs = append(errs, result.Error)
		} else {
			cards = append(cards, result.Cards)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return slices.Concat(cards...), nil
}

func fetchCollectionBatch(scryfallIds []string) ([]Card, error) {
	var cards []Card

	collectionUrl, err := url.JoinPath(scryfallUrl, "/cards/collection")
	if err != nil {
		return cards, err
	}

	query := collectionQuery{
		Identifiers: toScryfallIdentifiers(scryfallIds),
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

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return cards, fmt.Errorf("unexpected response content %s", contentType)
	}

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
