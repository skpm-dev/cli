package github

import (
	"fmt"
	"net/http"
	"os"
)

const baseURL = "https://api.github.com"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	token := os.Getenv("SKPM_GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("SKPM_GITHUB_TOKEN is not set — add it to your .zshrc or .bashrc")
	}

	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) newRequest(method, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req
}
