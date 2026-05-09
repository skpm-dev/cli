package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manifest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Version     string            `json:"version"`
	Repo        string            `json:"repo"`
	Skript      string            `json:"skript"`
	Minecraft   string            `json:"minecraft"`
	Addons      map[string]string `json:"addons"`
	Files       []string          `json:"files"`
}

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", path, err)
	}

	return &m, nil
}

func Validate(m *Manifest) error {
	required := map[string]string{
		"name":        m.Name,
		"description": m.Description,
		"author":      m.Author,
		"version":     m.Version,
		"repo":        m.Repo,
	}

	for field, value := range required {
		if value == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	if len(m.Files) == 0 {
		return fmt.Errorf("files must contain at least one .sk file")
	}

	return nil
}

func Save(m *Manifest, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}

	return nil
}
