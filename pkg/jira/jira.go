package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Client is a wrapper around the JIRA API
type Client struct {
	server   string
	username string
	password string
}

// New creates a new JIRA client
func New(server, username, password string) *Client {
	return &Client{server, username, password}
}

// Issue is a JIRA issue
type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary  string `json:"summary"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Status struct {
			Name string `json:"name"`
		}
	} `json:"fields"`
}

// Assignee returns the name of the assignee or 'Unassigned'
func (i *Issue) Assignee() string {
	if i.Fields.Assignee.DisplayName != "" {
		return i.Fields.Assignee.DisplayName
	}

	return "Unassigned"
}

// SearchResponse from JIRA
type SearchResponse struct {
	Issues     []Issue `json:"issues"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
}

// FindIssues finds JIRA issues matching a given JQL
func (c *Client) FindIssues(ctx context.Context, jql string) (*SearchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/rest/api/2/%s", c.server, "search"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("jql", jql)
	query.Set("fields", "key,summary,status,assignee")
	req.URL.RawQuery = query.Encode()

	req.SetBasicAuth(c.username, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response body: %s", body)
		return nil, fmt.Errorf("got unexpected status %d %s from JIRA", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	searchResp := SearchResponse{}
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// LinkForIssue returns a link to the given issue
func (c *Client) LinkForIssue(i *Issue) string {
	return fmt.Sprintf("https://%s/issues/%s", c.server, i.Key)
}
