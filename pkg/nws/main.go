package nws

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
	CountyCode string
	Office     string
	HTTPClient *http.Client
}

type errorResponse struct {
	Code    int    `json:"status"`
	Message string `json:"detail"`
}

func NewClient() *Client {
	apiURL := os.Getenv("NWS_API")
	if apiURL == "" {
		log.Fatal("NWS_API env var not found")
	}

	return &Client{
		BaseURL:    apiURL,
		CountyCode: os.Getenv("NWS_COUNTY"),
		Office:     os.Getenv("NWS_OFFICE"),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

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

	// Handle the response
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	fmt.Println(v)

	return nil
}
