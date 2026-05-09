package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/skpm-dev/cli/internal/manifest"
	"github.com/skpm-dev/cli/internal/models"
)

const registryURL = "https://skpm-registry-production.up.railway.app"

type publishRequest struct {
	Manifest manifest.Manifest `json:"manifest"`
	Files    map[string]string `json:"files"`
}

func GetPackage(name string) (*models.Package, error) {
	resp, err := http.Get(fmt.Sprintf("%s/packages/%s", registryURL, name))
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/publish", registryURL), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registry returned %d: %s", resp.StatusCode, string(b))
	}

	return nil
}
