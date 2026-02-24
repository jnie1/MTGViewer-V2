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
	"strconv"
	"strings"
)

var scryfallUrl = "https://api.scryfall.com"

func SearchCards(query string, page int) (SearchCardPage, error) {
	query, err := url.QueryUnescape(query)
	if err != nil {
		return SearchCardPage{}, err
	}

	searchPath, err := url.JoinPath(scryfallUrl, "/cards/search")
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

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return SearchCardPage{}, fmt.Errorf("unexpected response content %s", contentType)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchCardPage{}, err
	}

	var result searchResult
	if err := json.Unmarshal(content, &result); err != nil {
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

func FetchCard(scryfallId ScryfallIdentifier) (Card, error) {
	randomUrl, err := url.JoinPath(scryfallUrl, fmt.Sprintf("/cards/%s", scryfallId.Id))
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
	if result.CardFaces != nil && len(result.CardFaces) > 0 {
		result.Images = result.CardFaces[0].Images
	}
	return toCard(result), nil
}

func FetchCollection[Id CardIdentifier](identifiers []Id) ([]Card, error) {
	if len(identifiers) == 0 {
		return nil, errors.New("no ids are specified")
	}

	batchSizeLimit := 75

	results := make(chan collectionBatchResult)
	workerCount := 0

	for batch := range slices.Chunk(identifiers, batchSizeLimit) {
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

func fetchCollectionBatch[Id CardIdentifier](identifiers []Id) ([]Card, error) {
	collectionUrl, err := url.JoinPath(scryfallUrl, "/cards/collection")
	if err != nil {
		return nil, err
	}

	query := CollectionQuery[Id]{identifiers}
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
