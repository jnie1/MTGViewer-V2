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

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	q.Join("JOIN sets AS s ON s.code = c.setCode")
	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
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

	var rows []Card
	sql, params := q.Build()
	sql += " USING SAMPLE 1"

	if err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...); err != nil {
		return Card{}, err
	}

	if len(rows) == 0 {
		return Card{}, ErrMissingCards
	}

	row := rows[0]
	images, err := ImageURLs(row.ScryfallId)
	if err != nil {
		return Card{}, err
	}

	row.Images = images
	return row, nil
}

func FetchCard(ctx context.Context, scryfallId uuid.UUID) (Card, error) {
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return Card{}, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	q.Join("JOIN sets AS s ON s.code = c.setCode")
	q.WhereEq("ci.scryfallId", scryfallId)
	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
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
	q.Limit(1)

	var rows []Card
	sql, params := q.Build()

	if err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...); err != nil {
		return Card{}, err
	}

	if len(rows) == 0 {
		return Card{}, fmt.Errorf("no matching card found for %s", scryfallId)
	}

	row := rows[0]
	images, err := ImageURLs(row.ScryfallId)
	if err != nil {
		return Card{}, err
	}

	row.Images = images
	return row, nil
}

func FetchCollection(ctx context.Context, scryfallIds uuid.UUIDs) ([]Card, error) {
	if len(scryfallIds) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "sets"); err != nil {
		return nil, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	q.Join("JOIN sets AS s ON s.code = c.setCode")

	// TODO: maybe chunking?
	vals := make([]any, len(scryfallIds))
	for i, id := range scryfallIds {
		vals[i] = id
	}
	q.WhereIn("ci.scryfallId", vals)

	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
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

	var rows []Card
	sql, params := q.Build()

	if err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...); err != nil {
		return nil, err
	}

	for i := range rows {
		result := &rows[i]
		images, err := ImageURLs(result.ScryfallId)
		if err != nil {
			return nil, err
		}
		result.Images = images
	}

	return rows, nil
}

func FetchIdsByMultiverseId(ctx context.Context, multiverseIds []int) ([]CardId, error) {
	if len(multiverseIds) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")

	vals := make([]string, len(multiverseIds))
	for i, id := range multiverseIds {
		vals[i] = fmt.Sprintf("%d", id)
	}
	whereValues(q, "ci.multiverseId", vals)

	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
		"c.name",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
	)

	var rows []CardId
	sql, params := q.Build()
	err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...)

	return rows, err
}

func FetchIdsBySetCollector(ctx context.Context, setCollectors []SetCollectorNumber) ([]CardId, error) {
	if len(setCollectors) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")

	tups := make([][2]string, len(setCollectors))
	for i, sn := range setCollectors {
		tups[i] = [2]string{sn.Set, sn.CollectorNumber}
	}
	whereTuples2(q, [2]string{"c.setCode", "c.number"}, tups)

	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
		"c.name",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
	)

	var rows []CardId
	query, params := q.Build()
	err := sdk.Connection().ExecuteInto(ctx, &rows, query, params...)

	return rows, err
}

func FetchIdsByNameSet(ctx context.Context, nameSets []NameSet) ([]CardId, error) {
	if len(nameSets) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")

	tups := make([][2]string, len(nameSets))
	for i, ns := range nameSets {
		tups[i] = [2]string{ns.Name, ns.Set}
	}
	whereTuples2(q, [2]string{"c.name", "c.setCode"}, tups)

	q.Select(
		"ci.scryfallId",
		"ci.scryfallOracleId AS oracleId",
		"c.name",
		"c.setCode",
		"c.number AS collectorNumber",
		"CAST(ci.multiverseId AS INTEGER) AS multiverseId",
	)

	var rows []CardId
	sql, params := q.Build()
	err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...)

	return rows, err
}

func FetchScryfallIds(ctx context.Context, name string) ([]ScryfallIdObj, error) {
	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers"); err != nil {
		return nil, err
	}

	q := db.NewSQLBuilder("cards AS c")
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = c.uuid")
	q.WhereEq("c.name", name)
	q.Select("ci.scryfallId")

	var rows []ScryfallIdObj
	sql, params := q.Build()
	err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...)

	return rows, err
}

func FetchPrices(ctx context.Context, scryfallIds uuid.UUIDs, price float64) ([]CardPricePreview, error) {
	if len(scryfallIds) == 0 {
		return nil, nil
	}

	if err := sdk.EnsureViews(ctx, "cards", "card_identifiers", "all_prices_today"); err != nil {
		return nil, err
	}

	// TODO: check if needed
	// if _, err := sdk.Refresh(ctx); err != nil {
	// 	return nil, err
	// }

	q := priceDBQuery("p", price)
	q.Join("JOIN card_identifiers AS ci ON ci.uuid = p.uuid")

	// params should line up with outer query
	maxDate, _ := priceDBQuery("p2", price).Where("p2.uuid = ci.uuid").Select("MAX(p2.date)").Build()

	vals := make([]any, len(scryfallIds))
	for i, id := range scryfallIds {
		vals[i] = id
	}

	q.Where(fmt.Sprintf("p.date = (%s)", maxDate))
	q.WhereIn("ci.scryfallId", vals)
	q.Select("ci.scryfallId", "p.price")

	var rows []CardPricePreview
	sql, params := q.Build()
	err := sdk.Connection().ExecuteInto(ctx, &rows, sql, params...)

	return rows, err
}

func priceDBQuery(alias string, price float64) *db.SQLBuilder {
	q := db.NewSQLBuilder(fmt.Sprintf("all_prices_today AS %s", alias))

	q.Where(fmt.Sprintf("%s.provider = $1", alias), "tcgplayer")
	q.Where(fmt.Sprintf("%s.finish = $2", alias), "normal")

	q.Where(fmt.Sprintf("%s.price_type = $3", alias), "retail")
	q.Where(fmt.Sprintf("%s.price < $4", alias), price)

	return q
}
