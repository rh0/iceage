package toot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL    string
	authToken  string
	HTTPClient *http.Client
}

// Type to handle successful http responses
type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Type to handle http response errors
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Status type to handle a POST payload for a simple status
type Toot struct {
	Status string `json:"status"`
}

func NewClient(authToken string) *Client {
	apiURL := os.Getenv("MASTODON_API")
	if apiURL == "" {
		log.Fatal("MASTODON_API env var not found.")
	}

	return &Client{
		BaseURL:   apiURL,
		authToken: authToken,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) PostToot() error {
	postUrl := fmt.Sprintf("%s/statuses", c.BaseURL)

	buf := new(bytes.Buffer)
	t := Toot{
        Status: "I'm a bot: Beep Boop",
	}
	if err := json.NewEncoder(buf).Encode(&t); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postUrl, buf)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req); err != nil {
		return err
	}

	return nil
}

func (c *Client) sendRequest(req *http.Request) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fmt.Println(res.Body)

	// fullResponse := successResponse{
	// 	Data: v,
	// }
	// if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
	//     return err
	// }

	return nil
}
