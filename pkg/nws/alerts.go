package nws

import (
	"fmt"
	"net/http"
)

var endpoint = "alerts"

type AlertList struct {
	Features []Alert `json:"features"`
}

type Alert struct {
	IDURL      string          `json:"id"`
	Properties AlertProperties `json:"properties"`
}

type AlertProperties struct {
	ID          string `json:"id"`
	Sent        string `json:"sent"`
	Status      string `json:"status"`
	MessageType string `json:"messageType"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Certainty   string `json:"certainty"`
	Urgency     string `json:"urgency"`
	Event       string `json:"event"`
	SenderName  string `json:"senderName"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	Response    string `json:"response"`
}

func (c *Client) FetchAlerts() (AlertList, error) {
	var res AlertList
	getURL := fmt.Sprintf("%s/%s?zone=%s", c.BaseURL, endpoint, c.CountyCode)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return res, err
	}

	if err = c.sendRequest(req, &res); err != nil {
		return res, err
	}

	return res, nil
}
