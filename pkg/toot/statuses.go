package toot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"
)

// Status type to handle a POST payload for a simple status
type Toot struct {
	Spoiler   string `json:"spoiler_text"`
	Status    string `json:"status"`
	InReplyTo string `json:"in_reply_to_id"`
}

var endpoint = "statuses"

func (c *Client) Toot(message Toot) error {
	var toots []Toot

	// Check toot for length
	if err := c.validateToot(message); err != nil && err.Error() == "toot length is too long" {
		// The toot is too long, so split it up and thread it
		err := c.splitLongToot(message, &toots)
		if err != nil {
			return err
		}

		// Handle the toot thread
		err = c.multiToot(toots)
		if err != nil {
			return err
		}

		// We done
		return nil
	} else if err != nil {
		return err
	}

	// The toot was valid, so post it
	_, err := c.postToot(message)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) multiToot(toots []Toot) error {
	var lastID string
	var err error

	numToots := len(toots)

	for i, toot := range toots {
		pager := fmt.Sprintf(" (%d/%d)", i+1, numToots)

		toot.Spoiler = toot.Spoiler + pager
		toot.InReplyTo = lastID
		lastID, err = c.postToot(toot)
		if err != nil {
			return err
		}
	}

	return nil
}

// postToot Actually fire the POST for a toot.
func (c *Client) postToot(message Toot) (string, error) {
	postUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	var tootRes StatusResponse

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&message); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", postUrl, buf)
	if err != nil {
		return "", err
	}

	if tootRes, err = c.sendRequest(req); err != nil {
		return "", err
	}

	return tootRes.ID, nil
}

// validateToot Checks that the toot is of a valid form. Notably that it is
// less than 500 chars.
func (c *Client) validateToot(message Toot) error {
	if utf8.RuneCountInString(message.Spoiler)+utf8.RuneCountInString(message.Status) > 500 {
		return errors.New("toot length is too long")
	}
	return nil
}

// splitLongToot A recursive function to split up a long toot into multiple
// toots that will stay within the 500 char limit.
func (c *Client) splitLongToot(orig Toot, multiToot *[]Toot) error {
	remainingToot := Toot{
		Spoiler: orig.Spoiler,
		Status:  orig.Status,
	}
	spoilerLen := utf8.RuneCountInString(orig.Spoiler)

	firstStr := orig.Status[0 : 500-spoilerLen]
	strLen := utf8.RuneCountInString(firstStr)

	for i := strLen; i > 2; i-- {
		if firstStr[i-2:i] == "-=" {
			*multiToot = append(*multiToot, Toot{
				Spoiler: orig.Spoiler,
				Status:  firstStr[0 : i-2],
			})
			remainingToot.Status = orig.Status[i-2:]
			break
		}
	}

	if err := c.validateToot(remainingToot); err != nil && err.Error() == "toot length is too long" {
		// If what's left is still too long go through it all again
		c.splitLongToot(remainingToot, multiToot)
	} else if err != nil {
		return err
	} else {
		// We're good, tack on what's left
		*multiToot = append(*multiToot, remainingToot)
	}

	return nil
}
