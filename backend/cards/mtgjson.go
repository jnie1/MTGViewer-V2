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

type mtgJsonId struct {
	ScryfallId      uuid.UUID `json:"scryfallId"`
	Name            string    `json:"name"`
	Set             string    `json:"setCode"`
	CollectorNumber string    `json:"number"`
	MultiverseId    string    `json:"multiverseId,omitempty"`
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

	return fromJsonCard(results[0])
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

	return fromJsonCard(matches[0])
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
	var matches []mtgJsonCard

	if err := sdk.Connection().ExecuteInto(ctx, &matches, query, params...); err != nil {
		return nil, err
	}

	results := make([]Card, len(matches))

	for i, s := range matches {
		id, err := fromJsonCard(s)
		if err != nil {
			return nil, err
		}
		results[i] = id
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
		"ci.multiverseId",
		"c.name",
		"c.setCode",
		"c.number",
	)

	query, params := sql.Build()
	var matches []mtgJsonId

	if err := sdk.Connection().ExecuteInto(ctx, &matches, query, params...); err != nil {
		return nil, err
	}

	results := make([]CardId, len(matches))

	for i, m := range matches {
		id, err := fromJsonId(m)
		if err != nil {
			return nil, err
		}
		results[i] = id
	}

	return results, nil
}

func fromJsonCard(source mtgJsonCard) (Card, error) {
	multiverseId, err := parseMultiverseId(source.MultiverseId)
	if err != nil {
		return Card{}, err
	}

	images, err := ImageURLs(source.ScryfallId)
	if err != nil {
		return Card{}, err
	}

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

	return card, nil
}

func fromJsonId(source mtgJsonId) (CardId, error) {
	multiverseId, err := parseMultiverseId(source.MultiverseId)
	if err != nil {
		return CardId{}, err
	}

	card := CardId{
		source.ScryfallId,
		source.Name,
		source.Set,
		source.CollectorNumber,
		multiverseId,
	}

	return card, nil
}

func parseMultiverseId(source string) (int, error) {
	if source != "" {
		return strconv.Atoi(source)
	} else {
		return 0, nil
	}
}
