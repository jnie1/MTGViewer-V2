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
		return Card{}, ErrMissingCards
	}

	result := results[0]
	return fromMtgJson(result), nil
}

func TestFiltering(ctx context.Context) ([]Card, error) {
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return nil, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	sql.Join("JOIN sets AS s ON s.code = c.setCode")
	sql.Where("(c.setCode, c.number) IN (($1, $2), ($3, $4))", "DFT", "75", "TDM", "94")
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
	var results []mtgJsonCard

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return nil, err
	}

	return fromMtgJsons(results), nil
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
		return Card{}, fmt.Errorf("no matching card found for %s", scryfallId)
	}

	match := matches[0]
	return fromMtgJson(match), nil
}
