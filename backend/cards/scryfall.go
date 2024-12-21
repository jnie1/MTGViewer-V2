package cards

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ScryfallCard struct {
	ScryfallId    string `json:"scryfall_id"`
	MultiverseIds []int  `json:"multiverse_ids"`
	ManaCost      string `json:"mana_cost"`
	Name          string `json:"name"`
	Power         string `json:"power"`
	Toughness     string `json:"toughness"`
}

var scryfallUrl = "https://api.scryfall.com"

func FetchRandomCard() (ScryfallCard, error) {
	card := ScryfallCard{}
	randomUrl, err := url.JoinPath(scryfallUrl, "/cards/random")
	if err != nil {
		return card, err
	}

	req, err := http.NewRequest("GET", randomUrl, nil)
	if err != nil {
		return card, err
	}

	req.Header.Set("User-Agent", "mtg-viewer-v2")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return card, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return card, err
	}

	err = json.Unmarshal(body, &card)

	return card, err
}
