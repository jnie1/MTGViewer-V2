package containers

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

func ParseCardRequests(formFile *multipart.FileHeader) ([]CardRequest, error) {
	contentType := formFile.Header.Values("content-type")

	if slices.Contains(contentType, "text/plain") {
		return ParseTextFile(formFile)
	}

	if slices.Contains(contentType, "text/csv") {
		return ParseCsvFile(formFile)
	}

	return nil, errors.New("invalid file format")
}

func ParseTextFile(formFile *multipart.FileHeader) ([]CardRequest, error) {
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

		if !cardEntryPattern.MatchString(line) {
			return nil, fmt.Errorf("unexpected card format encountered: %s", line)
		}

		segments := cardEntryPattern.SubexpNames()

		name := segments[cardEntryPattern.SubexpIndex("name")]
		setCode := segments[cardEntryPattern.SubexpIndex("set")]
		collectorNumber := segments[cardEntryPattern.SubexpIndex("collector")]

		amount, err := strconv.Atoi(segments[cardEntryPattern.SubexpIndex("amount")])
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

	fetchedCards, err := cards.FetchCollection(setCollectors)
	if err != nil {
		return nil, err
	}

	requests := make([]CardRequest, len(fetchedCards))

	for i, card := range fetchedCards {
		source := cards.SetCollectorNumber{Set: card.SetCode, CollectorNumber: card.CollectorNumber}
		newRequest := CardRequest{ScryfallId: card.ScryfallId, Delta: amountMap[source]}
		requests[i] = newRequest
	}

	return requests, nil
}

func ParseCsvFile(formFile *multipart.FileHeader) ([]CardRequest, error) {
	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()
	csvReader := csv.NewReader(file)

	header, err := csvReader.Read()
	if err == io.EOF {
		return nil, errors.New("empty csv file received")
	}
	if err != nil {
		return nil, err
	}

	headerPositions := getHeaderPositions(header)
	if !headerPositions.hasValidPosition() {
		return nil, errors.New("invalid csv header, expected: scryfall id, multiverse id, or set code/collector number")
	}

	requests := []CardRequest{}

	multiverseIds := []cards.MultiverseIdentifier{}
	setCollectors := []cards.SetCollectorNumber{}
	nameSets := []cards.NameSet{}

	quantityMap := map[any]int{}

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

			newRequest := CardRequest{ScryfallId: scryfallId, Delta: quantity}
			requests = append(requests, newRequest)

		case headerPositions.MultiverseId > -1:
			multiverseId, err := strconv.Atoi(row[headerPositions.MultiverseId])
			if err != nil {
				return nil, err
			}

			newId := cards.MultiverseIdentifier{MultiverseId: multiverseId}
			multiverseIds = append(multiverseIds, newId)
			quantityMap[newId] = quantity

		case headerPositions.SetCode > -1 && headerPositions.CollectorNumber > -1:
			newId := cards.SetCollectorNumber{Set: row[headerPositions.SetCode], CollectorNumber: row[headerPositions.CollectorNumber]}
			setCollectors = append(setCollectors, newId)
			quantityMap[newId] = quantity

		case headerPositions.Name > -1 && headerPositions.SetCode > -1:
			newId := cards.NameSet{Name: row[headerPositions.Name], Set: row[headerPositions.SetCode]}
			nameSets = append(nameSets, newId)
			quantityMap[newId] = quantity
		}
	}

	if len(multiverseIds) > 0 {
		extraCards, err := cards.FetchCollection(multiverseIds)
		if err != nil {
			return nil, err
		}

		for _, card := range extraCards {
			if len(card.MultiverseIds) == 0 {
				return nil, fmt.Errorf("card resolved with no multiverse id: %s, (%s) %s", card.Name, card.SetCode, card.CollectorNumber)
			}
			source := cards.MultiverseIdentifier{MultiverseId: card.MultiverseIds[0]}
			requests = append(requests, CardRequest{ScryfallId: card.ScryfallId, Delta: quantityMap[source]})
		}
	}

	if len(setCollectors) > 0 {
		extraCards, err := cards.FetchCollection(setCollectors)
		if err != nil {
			return nil, err
		}

		for _, card := range extraCards {
			source := cards.SetCollectorNumber{Set: card.SetCode, CollectorNumber: card.CollectorNumber}
			requests = append(requests, CardRequest{ScryfallId: card.ScryfallId, Delta: quantityMap[source]})
		}
	}

	if len(nameSets) > 0 {
		extraCards, err := cards.FetchCollection(setCollectors)
		if err != nil {
			return nil, err
		}

		for _, card := range extraCards {
			source := cards.NameSet{Name: card.Name, Set: card.SetCode}
			requests = append(requests, CardRequest{ScryfallId: card.ScryfallId, Delta: quantityMap[source]})
		}
	}

	return requests, nil
}
