package cards

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	mtgjson "github.com/mtgjson/mtgjson-sdk-go"
	"github.com/mtgjson/mtgjson-sdk-go/db"
)

const (
	scryfallUrl     string = "https://api.scryfall.com"
	collectionLimit int    = 75
)

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

var missingCardsErr = errors.New("missing cards")

func FetchRandomCard() (Card, error) {
	sdk, err := mtgjson.New()
	if err != nil {
		return Card{}, err
	}

	defer sdk.Close()

	ctx := context.Background()
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return Card{}, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.Join("JOIN sets AS s ON s.code = c.setCode")
	sql.Select(
		"c.uuid",
		"ci.scryfallId",
		"c.manaCost",
		"c.name",
		"s.name AS setName",
		"c.setCode",
		"c.number",
		"ci.multiverseId",
		"c.power",
		"c.toughness",
		"c.type",
		"c.rarity",
	)

	query, params := sql.Build()
	query += " USING SAMPLE 1"

	var results []mtgJsonCard
	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return Card{}, err
	}

	if len(results) == 0 {
		return Card{}, missingCardsErr
	}

	result := results[0]
	return fromMtgJson(result), nil
}

func FetchCard(scryfallId ScryfallIdentifier) (Card, error) {
	sdk, err := mtgjson.New()
	if err != nil {
		return Card{}, err
	}

	defer sdk.Close()

	ctx := context.Background()
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return Card{}, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.Join("JOIN sets AS s ON s.code = c.setCode")
	sql.WhereEq("ci.scryfallId", scryfallId.Id)
	sql.Select(
		"c.uuid",
		"ci.scryfallId",
		"c.manaCost",
		"c.name",
		"s.name AS setName",
		"c.setCode",
		"c.number",
		"ci.multiverseId",
		"c.power",
		"c.toughness",
		"c.type",
		"c.rarity",
	)
	sql.Limit(1)

	query, params := sql.Build()
	var matches []mtgJsonCard

	if err := sdk.Connection().ExecuteInto(ctx, &matches, query, params...); err != nil {
		return Card{}, err
	}

	if len(matches) == 0 {
		return Card{}, fmt.Errorf("no matching card found for %s", scryfallId.Id)
	}

	match := matches[0]
	return fromMtgJson(match), nil
}

func FetchCollection[Id CardIdentifier](identifiers []Id) ([]Card, error) {
	if len(identifiers) == 0 {
		return nil, nil
	}

	results := make(chan collectionBatchResult)
	workerCount := 0

	for batch := range slices.Chunk(identifiers, collectionLimit) {
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
		if result.err != nil {
			errs = append(errs, result.err)
		} else {
			cards = append(cards, result.cards)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected response content %s", contentType)
	}

	decoder := json.NewDecoder(resp.Body)

	var result collectionResult
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return toCards(result.Cards), nil
}
