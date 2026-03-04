package containers

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/jnie1/MTGViewer-V2/cards"
)

var ErrFileFormat = errors.New("invalid file format")

func ParseCardRequests(formFile *multipart.FileHeader) ([]CardRequest, error) {
	fileExtension := filepath.Ext(formFile.Filename)

	if fileExtension == ".txt" {
		return parseTextFile(formFile)
	}

	if fileExtension == ".csv" {
		return parseCsvFile(formFile)
	}

	return nil, ErrFileFormat
}

func parseTextFile(formFile *multipart.FileHeader) ([]CardRequest, error) {
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

func parseCsvFile(formFile *multipart.FileHeader) ([]CardRequest, error) {
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

	requests := []CardRequest{}
	identifierOptions := map[int]cards.CardIdentifier{}
	extraIds := []cards.CardIdentifier{}
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
			identifierOptions[1] = newId
			extraIds = append(extraIds, newId)
			quantityMap[newId] = quantity

		case headerPositions.SetCode > -1 && headerPositions.CollectorNumber > -1:
			newId := cards.SetCollectorNumber{Set: row[headerPositions.SetCode], CollectorNumber: row[headerPositions.CollectorNumber]}
			identifierOptions[2] = newId
			extraIds = append(extraIds, newId)
			quantityMap[newId] = quantity

		case headerPositions.Name > -1 && headerPositions.SetCode > -1:
			newId := cards.NameSet{Name: row[headerPositions.Name], Set: row[headerPositions.SetCode]}
			identifierOptions[3] = newId
			extraIds = append(extraIds, newId)
			quantityMap[newId] = quantity
		}
	}

	if len(extraIds) > 0 {
		extraCards, err := cards.FetchCollection(extraIds)
		if err != nil {
			return nil, err
		}
		for _, card := range extraCards {
			for _, id := range identifierOptions {
				source, err := id.Convert(card)
				if err != nil {
					continue
				}
				if amount, ok := quantityMap[source]; ok {
					requests = append(requests, CardRequest{card.ScryfallId, amount})
					break
				}
			}
		}
	}

	return requests, nil
}
