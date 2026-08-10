package cards

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	mtgjson "github.com/mtgjson/mtgjson-sdk-go"
	"github.com/mtgjson/mtgjson-sdk-go/db"
)

var sdk *mtgjson.SDK
var ErrMissingCards = errors.New("missing cards")

type mtgJsonCard struct {
	MtgJsonId       uuid.UUID `json:"uuid"`
	ScryfallId      uuid.UUID `json:"scryfallId"`
	ManaCost        string    `json:"manaCost,omitempty"`
	Name            string    `json:"name"`
	SetName         string    `json:"setName"`
	Set             string    `json:"setCode"`
	CollectorNumber string    `json:"number"`
	MultiverseId    string    `json:"multiverseId,omitempty"`
	Power           string    `json:"power,omitempty"`
	Toughness       string    `json:"toughness,omitempty"`
	Type            string    `json:"type"`
	Rarity          string    `json:"rarity"`
}

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

func FetchIdentifiers(ids CardIdQuery) ([]CardId, error) {
	if ids.IsEmpty() {
		return nil, nil
	}

	ctx := context.Background()
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	sql := db.NewSQLBuilder("cards AS c")
	sql.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	WhereIds(sql, ids, "ci.multiverseId", "c.name", "c.setCode", "c.number")
	sql.Select(
		"ci.scryfallId",
		"ci.multiverseId",
		"c.name",
		"c.setCode",
		"c.number AS collectorNumber",
	)

	query, params := sql.Build()
	var results []CardId

	if err := sdk.Connection().ExecuteInto(ctx, &results, query, params...); err != nil {
		return nil, err
	}

	return results, nil
}

func fromMtgJson(source mtgJsonCard) Card {
	var multiverseId int
	if source.MultiverseId != "" {
		if id, err := strconv.Atoi(source.MultiverseId); err == nil {
			multiverseId = id
		}
	}

	images := CardImageUrls{}

	card := Card{
		source.ScryfallId,
		source.Name,
		source.ManaCost,
		source.SetName,
		source.Set,
		source.CollectorNumber,
		multiverseId,
		source.Type,
		source.Rarity,
		source.Power,
		source.Toughness,
		images,
		map[string]string{},
	}

	return card
}

func fromMtgJsons(sources []mtgJsonCard) []Card {
	results := make([]Card, len(sources))
	for i, s := range sources {
		results[i] = fromMtgJson(s)
	}
	return results
}
