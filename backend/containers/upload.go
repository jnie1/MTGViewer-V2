package containers

import (
	"bufio"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"slices"
	"strconv"

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
		sourceSetCollector := cards.SetCollectorNumber{Set: card.SetCode, CollectorNumber: card.CollectorNumber}
		delta := amountMap[sourceSetCollector]

		requests[i] = CardRequest{ScryfallId: card.ScryfallId, Delta: delta}
	}

	return requests, nil
}

func ParseCsvFile(formFile *multipart.FileHeader) ([]CardRequest, error) {
	file, err := formFile.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return nil, nil
}
