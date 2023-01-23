package toot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var endpoint = "statuses"

func (c *Client) PostToot(tootMessage string) error {
	postUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	buf := new(bytes.Buffer)
	t := Toot{
		Status: tootMessage,
	}
	if err := json.NewEncoder(buf).Encode(&t); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postUrl, buf)
	if err != nil {
		return err
	}

	if err = c.sendRequest(req); err != nil {
		return err
	}

	return nil
}
