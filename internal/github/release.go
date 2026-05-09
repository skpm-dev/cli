package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type Release struct {
	ID        int    `json:"id"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
}

func (c *Client) CreateRelease(owner, repo, tag, version string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", baseURL, owner, repo)

	body := map[string]string{
		"tag_name": tag,
		"name":     "v" + version,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github returned %d: %s", resp.StatusCode, string(b))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release response: %w", err)
	}

	return &release, nil
}

func (c *Client) UploadReleaseAsset(release *Release, filePath string) (string, error) {
	fileName := filepath.Base(filePath)
	uploadURL := cleanUploadURL(release.UploadURL, fileName)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read file %s: %w", filePath, err)
	}

	req, _ := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(fileData))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", mime.TypeByExtension(filepath.Ext(fileName)))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github returned %d: %s", resp.StatusCode, string(b))
	}

	var asset struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&asset); err != nil {
		return "", fmt.Errorf("failed to decode asset response: %w", err)
	}

	return asset.BrowserDownloadURL, nil
}

// GitHub's upload URL comes with a template suffix like {?name,label} — strip it
func cleanUploadURL(raw, fileName string) string {
	base := raw
	if idx := len(raw) - len("{?name,label}"); idx > 0 {
		base = raw[:idx]
	}
	return base + "?name=" + fileName
}
