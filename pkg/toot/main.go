package toot

import (
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

func (c *Client) sendRequest(req *http.Request) error {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	// Fire the request
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

	return nil
}
