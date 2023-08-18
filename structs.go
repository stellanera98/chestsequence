package main

import "fmt"

type AboutV2 struct {
	Gacha GachaV2 `json:"gacha"`
}

type GachaV2 struct {
	PlayerData interface{} `json:"player_data"`
	Params     GachaParams `json:"params"`
}

type GachaParams struct {
	DeckIndices map[string]int   `json:"deck_indices"`
	Decks       map[string][]int `json:"decks"`
}

type EventV2 struct {
	GachaV2 struct {
		Params struct {
			Decks map[string][]int `json:"decks"`
		} `json:"params"`
	} `json:"gacha"`
}

type DragonsongGacha struct {
	EventGacha map[string]Gacha `json:"gacha"`
}

type Gacha struct {
	DropLists map[string][]DropList `json:"drops"`
	SpinTypes []SpinType            `json:"spin_types"`
}

type DropList struct {
	Kind     string  `json:"kind"`
	Weight   int     `json:"weight"`
	DropType string  `json:"drop_type"`
	ID       string  `json:"id"`
	Amount   float64 `json:"mu"`
}

func (d DropList) String() string {
	return fmt.Sprintf("%s:%.0f", d.ID, d.Amount)
}

type SpinType struct {
	Title     string     `json:"title"`
	DropRates []DropRate `json:"drop_rates"`
}

type DropRate struct {
	Count  float64 `json:"count"`
	Chance float64 `json:"chance"`
	Rarity string  `json:"drop_type"`
	DropID string  `json:"id"`
}
