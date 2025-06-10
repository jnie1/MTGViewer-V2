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

	"github.com/google/uuid"
)

var scryfallUrl = "https://api.scryfall.com"

func FetchRandomCard() (Card, error) {
	randomUrl, err := url.JoinPath(scryfallUrl, "/cards/random")
	if err != nil {
		return Card{}, err
	}

	req, err := http.NewRequest("GET", randomUrl, nil)
	if err != nil {
		return Card{}, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Card{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return Card{}, fmt.Errorf("unexpected response content %s", contentType)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return Card{}, err
	}

	var result scryfallCard
	if err := json.Unmarshal(content, &result); err != nil {
		return Card{}, err
	}

	return toCard(result), nil
}

func FetchCard(scryfallId uuid.UUID) (Card, error) {
	randomUrl, err := url.JoinPath(scryfallUrl, fmt.Sprintf("/cards/%s", scryfallId))
	if err != nil {
		return Card{}, err
	}

	req, err := http.NewRequest("GET", randomUrl, nil)
	if err != nil {
		return Card{}, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Card{}, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return Card{}, fmt.Errorf("unexpected response content %s", contentType)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return Card{}, err
	}

	var result scryfallCard
	if err := json.Unmarshal(content, &result); err != nil {
		return Card{}, err
	}

	return toCard(result), nil
}

func FetchCollection(scryfallIds uuid.UUIDs) ([]Card, error) {
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

	for range workerCount {
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

func fetchCollectionBatch(scryfallIds uuid.UUIDs) ([]Card, error) {
	collectionUrl, err := url.JoinPath(scryfallUrl, "/cards/collection")
	if err != nil {
		return nil, err
	}

	query := collectionQuery{
		Identifiers: toScryfallIdentifiers(scryfallIds),
	}

	payload, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", collectionUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected response content %s", contentType)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result collectionResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return toCards(result.Cards), nil
}
