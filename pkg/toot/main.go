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

// Type to handle a successful responce
type StatusResponse struct {
	ID string `json:"id"`
}

// Type to handle http response errors
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

func (c *Client) sendRequest(req *http.Request) (StatusResponse, error) {
	var tootRes StatusResponse

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	// Fire the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return tootRes, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return tootRes, errors.New(errRes.Message)
		}

		return tootRes, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(&tootRes); err != nil {
		return tootRes, err
	}

	return tootRes, nil
}
