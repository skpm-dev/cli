package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	Login string `json:"login"`
}

func Token() (string, error) {
	token := os.Getenv("SKPM_GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("SKPM_GITHUB_TOKEN is not set — add it to your .zshrc or .bashrc")
	}
	return token, nil
}

func GetAuthenticatedUser(token string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not reach github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid github token — check SKPM_GITHUB_TOKEN")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github returned %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("could not decode github response: %w", err)
	}

	return &user, nil
}
