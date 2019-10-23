package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// BlockTypeSection is a section block which contains more content
	BlockTypeSection = "section"

	// BlockTypeDivider is a horizontal divider
	BlockTypeDivider = "divider"

	// BlockTextTypeMarkdown is the type for markdown formatted text
	BlockTextTypeMarkdown = "mrkdwn"
)

// Text is some structured text
type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Block represents block in Slack's Block Kit https://api.slack.com/block-kit
type Block struct {
	Type string `json:"type"`
	Text *Text  `json:"text,omitempty"`
}

// Message is a Slack message
type Message struct {
	Blocks []*Block `json:"blocks"`
}

// Client is a wrapper around the Slack API
type Client struct {
	webhookURL string
}

// New constructs a new slack.Client
func New(webhookURL string) *Client {
	return &Client{webhookURL}
}

// PostMessage posts a message to the client's webhook URL
func (c *Client) PostMessage(ctx context.Context, msg *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(encodedMsg))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("Non-200 response body: %s", body)
		return fmt.Errorf("got unexpected status %d %s from Slack", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
