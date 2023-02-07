package toot

import (
	// "bytes"
	// "encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"
	// "net/http"
)

// Status type to handle a POST payload for a simple status
type Toot struct {
	Spoiler string `json:"spoiler_text"`
	Status  string `json:"status"`
}

var endpoint = "statuses"

func (c *Client) PostToot(message Toot) error {
	// postUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	// Check that for length.
	if err := c.validateToot(message); err.Error() == "toot length is too long" {
		var multiToot []Toot
		err := c.splitLongToot(message, &multiToot)
		if err != nil {
			return err
		}

		fmt.Println(multiToot)
	}

	// buf := new(bytes.Buffer)
	// if err := json.NewEncoder(buf).Encode(&message); err != nil {
	// 	return err
	// }

	// req, err := http.NewRequest("POST", postUrl, buf)
	// if err != nil {
	// 	return err
	// }

	// if err = c.sendRequest(req); err != nil {
	// 	return err
	// }

	return nil
}

// validateToot Checks that the toot is of a valid form. Notably that is
// less than 500 chars.
func (c *Client) validateToot(message Toot) error {
	if utf8.RuneCountInString(message.Spoiler)+utf8.RuneCountInString(message.Status) > 500 {
		return errors.New("toot length is too long")
	}
	return nil
}

func (c *Client) splitLongToot(orig Toot, multiToot *[]Toot) error {
    remainingToot := Toot{
        Spoiler: orig.Spoiler,
        Status:  orig.Status,
    }
	spoilerLen := utf8.RuneCountInString(orig.Spoiler)

	firstStr := orig.Status[0 : 500-spoilerLen]
	strLen := utf8.RuneCountInString(firstStr)

	for i := strLen; i > 1; i-- {
		if firstStr[i-2:i] == "-=" {
            fmt.Println(firstStr[0 : i-2])
			*multiToot = append(*multiToot, Toot{
				Spoiler: orig.Spoiler,
				Status:  firstStr[0 : i-2],
			})
			remainingToot.Status = orig.Status[i-2:]
			break
		}
	}

    fmt.Println("REMAINING")
    fmt.Println(remainingToot)
    // If what's left is still too long go through it all again
    if err := c.validateToot(remainingToot); err.Error() == "toot length is too long" {
        fmt.Println("RECURSE")
        c.splitLongToot(remainingToot, multiToot)
    }

    // We're good, tack on what's left
	*multiToot = append(*multiToot, remainingToot)

	return nil
}
