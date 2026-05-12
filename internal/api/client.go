package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/skpm-dev/cli/internal/manifest"
	"github.com/skpm-dev/cli/internal/models"
)

func registryURL() string {
	if u := os.Getenv("SKPM_REGISTRY_URL"); u != "" {
		return u
	}
	return "https://skpm-registry-production.up.railway.app"
}

type publishRequest struct {
	Manifest manifest.Manifest `json:"manifest"`
	Files    map[string]string `json:"files"`
}

func GetPackage(name string) (*models.Package, error) {
	resp, err := http.Get(fmt.Sprintf("%s/packages/%s", registryURL(), name))
	if err != nil {
		return nil, fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	var pkg models.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return &pkg, nil
}

func Publish(token string, m *manifest.Manifest, files map[string]string) error {
	body := publishRequest{
		Manifest: *m,
		Files:    files,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not encode request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/publish", registryURL()), bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrVersionConflict
	}

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registry returned %d: %s", resp.StatusCode, string(b))
	}

	return nil
}

// Search queries the registry for packages matching the given query string.
func Search(query string) ([]models.PackageSummary, error) {
	u := fmt.Sprintf("%s/search?q=%s", registryURL(), url.QueryEscape(query))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	var results []models.PackageSummary
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}
	return results, nil
}

// ErrVersionConflict is returned when the version being published already has an open PR.
var ErrVersionConflict = fmt.Errorf("version already has an open pull request")

// Yank marks a specific version as yanked in the registry.
func Yank(adminToken, name, version, reason string) error {
	url := fmt.Sprintf("%s/packages/%s/%s", registryURL(), name, version)
	return adminDelete(adminToken, url, reason)
}

// Remove hard-removes an entire package from the registry.
func Remove(adminToken, name, reason string) error {
	url := fmt.Sprintf("%s/packages/%s", registryURL(), name)
	return adminDelete(adminToken, url, reason)
}

func adminDelete(adminToken, url, reason string) error {
	body := struct {
		Reason string `json:"reason"`
	}{Reason: reason}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid admin token")
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registry returned %d: %s", resp.StatusCode, b)
	}
	return nil
}
