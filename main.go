package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	WDVer   string //= "820"
	BaseURL = "https://%s-dot-pgdragonsong.appspot.com/"
)

var (
	WDUserAgentVersion string //= "8.20"
	WDUserAgentAvatar  string //= "820001"
	WDUserAgentBase    = "War Dragons/%s pgid=%s tag=1cfc6a1b83dcb8c4 cv=%s hw=SONY XPERIA XZ; kagura_dsds; gzip"
)

func main() {
	version := flag.String("version", "8.20", "should be adjusted for every new release of WD")
	flag.Parse()

	WDVer = strings.Replace(*version, ".", "", 1)
	WDUserAgentAvatar = fmt.Sprintf("%s001", WDVer)
	WDUserAgentVersion = *version

	var email string
	var passwd string

	fmt.Print("Pocket ID email: ")
	fmt.Scanln(&email)
	fmt.Print("Pocket ID password: ")
	fmt.Scanln(&passwd)
	pgid, pgid_token, err := getPGIDandToken(email, passwd)
	if err != nil {
		fmt.Println("couldnt get pgid and pgid_token:", err.Error())
		return
	}
	UpdateData(pgid, pgid_token)
}

func getPGIDandToken(email, passwd string) (string, string, error) {
	resp, err := http.Get(fmt.Sprintf("https://pocket-gems.appspot.com/pocket_id/login?password=%s&email=%s", passwd, email))
	if err != nil {
		return "", "", err
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var res struct {
		Pgid_token string `json:"pocket_id_token"`
		Pgid       string `json:"pocket_id"`
	}

	err = json.Unmarshal(out, &res)
	if err != nil {
		return "", "", err
	}

	return res.Pgid, res.Pgid_token, nil
}

func UpdateData(pgid, pgid_token string) {
	client := NewClient(pgid, pgid_token)
	clientData, err := client.GetEventData()
	if err != nil {
		fmt.Println(err)
		return
	}

	clientGachaData, err := client.GetGachaData()
	if err != nil {
		fmt.Println("error getting gacha data:", err.Error())
		return
	}

	name := strings.Builder{}
	b := false
	for k := range clientGachaData {
		if b {
			name.WriteRune('_')
		}
		name.WriteString(k)
		b = true
	}

	indexMap := make(map[string]map[string]int)
	for event, data := range clientGachaData {
		indexMap[event] = data.Gacha.Params.DeckIndices
	}

	SaveJson(fmt.Sprintf("client_%s.json", name.String()), indexMap)

	simple := make(map[string]map[string][]string)
	for eventID, gachaData := range clientGachaData {
		simple[eventID] = make(map[string][]string)
		for sequenceName, sequence := range gachaData.Gacha.Params.Decks {
			simple[eventID][sequenceName] = make([]string, len(sequence))
			for position, drop := range sequence {
				simple[eventID][sequenceName][position] = clientData.EventGacha[eventID].DropLists[sequenceName][drop].String()
			}
		}
	}

	SaveJson(fmt.Sprintf("%s.json", name.String()), simple)
}
