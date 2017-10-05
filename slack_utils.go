package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const slackAPI = "https://api.slack.com/"

type rtmStartResponse struct {
	Ok    bool         `json:"ok"`
	URL   string       `json:"url"`
	Error string       `json:"error"`
	Self  rtmStartSelf `json:"self"`
}

type rtmStartSelf struct {
	ID string `json:"id"`
}

type Event struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
	ID      uint64 `json:"id"`
}

func openSlackWebSocket(token string) (webSocketURL, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", response.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}
	var startResponse rtmStartResponse
	err = json.Unmarshal(body, &startResponse)
	if err != nil {
		return
	}

	if !startResponse.Ok {
		err = fmt.Errorf("Slack error: %s", startResponse.Error)
		return
	}

	webSocketURL = startResponse.URL
	id = startResponse.Self.ID
	return
}
