package cards

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	mtgjson "github.com/mtgjson/mtgjson-sdk-go"
	"github.com/mtgjson/mtgjson-sdk-go/db"
)

var sdk *mtgjson.SDK
var ErrMissingCards = errors.New("missing cards")

func OpenSDK() (*mtgjson.SDK, error) {
	var err error
	sdk, err = mtgjson.New()
	return sdk, err
}

func FetchRandomCard(ctx context.Context) (Card, error) {
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
		"s.name AS set",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
		"c.power",
		"c.toughness",
		"c.type",
		"c.rarity",
	)

	query, params := sql.Build()
	query += " USING SAMPLE 1"
	var results []Card

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return Card{}, err
	}

	if len(results) == 0 {
		return Card{}, ErrMissingCards
	}

	result := results[0]
	images, err := ImageURLs(result.ScryfallId)
	if err != nil {
		return Card{}, err
	}

	result.Images = images
	return result, nil
}

func FetchCard(ctx context.Context, scryfallId uuid.UUID) (Card, error) {
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return Card{}, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.Join("JOIN sets AS s ON s.code = c.setCode")
	sql.WhereEq("ci.scryfallId", scryfallId)
	sql.Select(
		"ci.scryfallId",
		"c.manaCost",
		"c.name",
		"s.name AS set",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
		"c.power",
		"c.toughness",
		"c.type",
		"c.rarity",
	)
	sql.Limit(1)

	query, params := sql.Build()
	var results []Card

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return Card{}, err
	}

	if len(results) == 0 {
		return Card{}, fmt.Errorf("no matching card found for %s", scryfallId)
	}

	result := results[0]
	images, err := ImageURLs(result.ScryfallId)
	if err != nil {
		return Card{}, err
	}

	result.Images = images
	return result, nil
}

func FetchCollection(ctx context.Context, scryfallIds uuid.UUIDs) ([]Card, error) {
	if len(scryfallIds) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return nil, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.Join("JOIN sets AS s ON s.code = c.setCode")

	// TODO: maybe chunking?
	vals := make([]any, len(scryfallIds))
	for i, id := range scryfallIds {
		vals[i] = id
	}
	sql.WhereIn("ci.scryfallId", vals)

	sql.Select(
		"ci.scryfallId",
		"c.manaCost",
		"c.name",
		"s.name AS set",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
		"c.power",
		"c.toughness",
		"c.type",
		"c.rarity",
	)

	query, params := sql.Build()
	var results []Card

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return nil, err
	}

	for i := range results {
		result := &results[i]
		images, err := ImageURLs(result.ScryfallId)
		if err != nil {
			return nil, err
		}
		result.Images = images
	}

	return results, nil

}

func FetchIdentifiers(ctx context.Context, ids CardIdQuery) ([]CardId, error) {
	if ids.IsEmpty() {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	WhereIds(sql, ids, "ci.multiverseId", "c.name", "c.setCode", "c.number")
	sql.Select(
		"ci.scryfallId",
		"c.name",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
	)

	query, params := sql.Build()
	var results []CardId

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return nil, err
	}

	return results, nil
}

func FetchScryfallIds(ctx context.Context, name string) ([]ScryfallIdObj, error) {
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.WhereEq("c.name", name)
	sql.Select("ci.scryfallId")

	query, params := sql.Build()
	var results []ScryfallIdObj

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return nil, err
	}

	return results, nil
}
