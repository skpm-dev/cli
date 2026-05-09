package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PullRequest struct {
	HTMLURL string `json:"html_url"`
}

func (c *Client) OpenPR(owner, repo, branch, title, body string) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", baseURL, owner, repo)

	payload := map[string]string{
		"title": title,
		"body":  body,
		"head":  branch,
		"base":  "main",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to open PR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github returned %d: %s", resp.StatusCode, string(b))
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode PR response: %w", err)
	}

	return &pr, nil
}
