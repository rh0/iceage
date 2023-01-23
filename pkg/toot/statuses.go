package toot

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

var endpoint = "statuses"

func (c *Client) PostToot(tootSpoiler, tootMessage string) error {
	postUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

    formData := url.Values{
        "status": {tootMessage[:400]},
        "spoiler_text": {tootSpoiler},
    }
    fmt.Println(formData)

	req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil { return err
	}

	if err = c.sendRequest(req); err != nil {
		return err
	}

	return nil
}
