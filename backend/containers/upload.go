package containers

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

var ErrFileFormat = errors.New("invalid file format")

func ParseCardRequests(ctx context.Context, formFile *multipart.FileHeader) ([]CardRequest, error) {
	fileExtension := filepath.Ext(formFile.Filename)

	if fileExtension == ".txt" {
		return parseTextFile(ctx, formFile)
	}

	if fileExtension == ".csv" {
		return parseCsvFile(ctx, formFile)
	}

	return nil, ErrFileFormat
}

func parseTextFile(ctx context.Context, formFile *multipart.FileHeader) ([]CardRequest, error) {
	cardEntryPattern, err := regexp.Compile(`^(?P<amount>\d+) (?P<name>.+?) \((?P<set>.+?)\) (?P<collector>.+)$`)
	if err != nil {
		return nil, err
	}

	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	setCollectors := []cards.SetCollectorNumber{}
	amountMap := map[cards.SetCollectorNumber]int{}

	for scanner.Scan() {
		line := scanner.Text()

		match := cardEntryPattern.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("unexpected card format encountered: %s", line)
		}

		name := match[cardEntryPattern.SubexpIndex("name")]
		setCode := match[cardEntryPattern.SubexpIndex("set")]
		collectorNumber := match[cardEntryPattern.SubexpIndex("collector")]

		amount, err := strconv.Atoi(match[cardEntryPattern.SubexpIndex("amount")])
		if err != nil {
			return nil, err
		}
		if amount <= 0 {
			return nil, fmt.Errorf("invalid amount for %s: %d", name, amount)
		}

		newEntry := cards.SetCollectorNumber{Set: setCode, CollectorNumber: collectorNumber}
		setCollectors = append(setCollectors, newEntry)
		amountMap[newEntry] = amount
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	cardIds, err := cards.FetchIdsBySetCollector(ctx, setCollectors)
	if err != nil {
		return nil, err
	}

	requests := make([]CardRequest, len(cardIds))

	for i, card := range cardIds {
		source := cards.SetCollectorNumber{Set: card.SetCode, CollectorNumber: card.CollectorNumber}
		newRequest := CardRequest{card.ScryfallId, card.OracleId, amountMap[source]}
		requests[i] = newRequest
	}

	return requests, nil
}

func parseCsvFile(ctx context.Context, formFile *multipart.FileHeader) ([]CardRequest, error) {
	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()
	csvReader := csv.NewReader(file)

	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	headerPositions := getHeaderPositions(header)
	if !headerPositions.Valid() {
		return nil, csv.ErrFieldCount
	}

	scryfallIds := map[uuid.UUID]int{}
	multiverseIds := map[int]int{}
	setCollectors := map[cards.SetCollectorNumber]int{}
	nameSets := map[cards.NameSet]int{}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(row[headerPositions.Quantity])
		if err != nil {
			return nil, err
		}

		switch {
		case headerPositions.ScryfallId > -1:
			scryfallId, err := uuid.Parse(row[headerPositions.ScryfallId])
			if err != nil {
				return nil, err
			}
			scryfallIds[scryfallId] = quantity

		case headerPositions.MultiverseId > -1:
			multiverseId, err := strconv.Atoi(row[headerPositions.MultiverseId])
			if err != nil {
				return nil, err
			}
			multiverseIds[multiverseId] = multiverseIds[multiverseId] + quantity

		case headerPositions.SetCode > -1 && headerPositions.CollectorNumber > -1:
			setCollector := cards.SetCollectorNumber{
				Set:             row[headerPositions.SetCode],
				CollectorNumber: row[headerPositions.CollectorNumber],
			}
			setCollectors[setCollector] = setCollectors[setCollector] + quantity

		case headerPositions.Name > -1 && headerPositions.SetCode > -1:
			nameSet := cards.NameSet{
				Name: row[headerPositions.Name],
				Set:  row[headerPositions.SetCode],
			}
			nameSets[nameSet] = nameSets[nameSet] + quantity
		}
	}

	requests := []CardRequest{}

	if len(scryfallIds) > 0 {
		keys := slices.Collect(maps.Keys(scryfallIds))
		fullCards, err := cards.FetchCollection(ctx, keys)
		if err != nil {
			return nil, err
		}
		for _, card := range fullCards {
			if amount, ok := scryfallIds[card.ScryfallId]; ok {
				requests = append(requests, CardRequest{card.ScryfallId, card.OracleId, amount})
			}
		}
	}

	if len(multiverseIds) > 0 {
		keys := slices.Collect(maps.Keys(multiverseIds))
		cardIds, err := cards.FetchIdsByMultiverseId(ctx, keys)
		if err != nil {
			return nil, err
		}
		for _, card := range cardIds {
			if amount, ok := multiverseIds[card.MultiverseId]; ok {
				requests = append(requests, CardRequest{card.ScryfallId, card.OracleId, amount})
			}
		}
	}

	if len(setCollectors) > 0 {
		keys := slices.Collect(maps.Keys(setCollectors))
		cardIds, err := cards.FetchIdsBySetCollector(ctx, keys)
		if err != nil {
			return nil, err
		}
		for _, card := range cardIds {
			if amount, ok := setCollectors[card.SetCollectorNumber()]; ok {
				requests = append(requests, CardRequest{card.ScryfallId, card.OracleId, amount})
			}
		}
	}

	if len(nameSets) > 0 {
		keys := slices.Collect(maps.Keys(nameSets))
		cardIds, err := cards.FetchIdsByNameSet(ctx, keys)
		if err != nil {
			return nil, err
		}
		for _, card := range cardIds {
			if amount, ok := nameSets[card.NameSet()]; ok {
				requests = append(requests, CardRequest{card.ScryfallId, card.OracleId, amount})
			}
		}
	}

	return requests, nil
}
