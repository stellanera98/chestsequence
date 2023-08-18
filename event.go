package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Client struct {
	Client        *http.Client
	PocketID      string
	PocketIDToken string
}

func NewClient(id, idToken string) *Client {
	c := &Client{
		PocketID:      id,
		PocketIDToken: idToken,
		Client:        &http.Client{},
	}
	j, _ := cookiejar.New(nil)
	c.Client.Jar = j

	return c
}

func (c *Client) eventDataURL() string {
	return fmt.Sprintf("%s/dragons/event/current?player_id=%s&pocket_id=%s&pocket_id_token=%s&tab=&native_gacha_enabled=0&season_tutorial_enabled=0", fmt.Sprintf(BaseURL, WDVer), c.PocketID, c.PocketID, c.PocketIDToken)
}

func (c *Client) gachaDataURL() string {
	return fmt.Sprintf("%s/ext/dragonsong/event/about_v2", fmt.Sprintf(BaseURL, WDVer))
}

// GetEventData just gives the gacha part of it because thats what I need here. TODO: fix
func (c *Client) GetEventData() (DragonsongGacha, error) {
	var ret DragonsongGacha
	request, err := http.NewRequest(http.MethodGet, c.eventDataURL(), nil)
	if err != nil {
		return DragonsongGacha{}, err
	}

	request.Header.Add("user-agent", fmt.Sprintf(WDUserAgentBase, WDUserAgentVersion, WDUserAgentAvatar, c.PocketID))
	request.Header.Add("Accept-Encoding", "gzip")
	resp, err := c.Client.Do(request)
	if err != nil {
		return DragonsongGacha{}, nil
	}

	tmp, err := gzip.NewReader(resp.Body)
	if err != nil {
		return DragonsongGacha{}, err
	}

	out, err := io.ReadAll(tmp)
	if err != nil {
		return DragonsongGacha{}, nil
	}

	start := strings.Index(string(out), "window.params_and_data = ") + 25
	end := strings.Index(string(out), "}};") + 2
	if start-25 == -1 || end-2 == -1 {
		return DragonsongGacha{}, fmt.Errorf(string(out))
	}

	err = json.Unmarshal(out[start:end], &ret)
	return ret, err
}

func (c *Client) GetGachaData() (map[string]AboutV2, error) {
	var ret map[string]AboutV2
	request, err := http.NewRequest(http.MethodGet, c.gachaDataURL(), nil)
	request.Header.Add("user-agent", fmt.Sprintf(WDUserAgentBase, WDUserAgentVersion, WDUserAgentAvatar, c.PocketID))
	request.Header.Add("accept-encoding", "gzip")
	if err != nil {
		return map[string]AboutV2{}, err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return map[string]AboutV2{}, nil
	}

	tmp, err := gzip.NewReader(resp.Body)
	if err != nil {
		return map[string]AboutV2{}, err
	}

	out, err := io.ReadAll(tmp)
	if err != nil {
		return map[string]AboutV2{}, nil
	}

	err = json.Unmarshal(out, &ret)
	return ret, err
}
