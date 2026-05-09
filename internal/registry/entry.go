package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/skpm-dev/cli/internal/manifest"
)

const registryBase = "https://raw.githubusercontent.com/skpm-dev/registry/main"

type FileEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionEntry struct {
	Skript    string            `json:"skript"`
	Minecraft string            `json:"minecraft"`
	Addons    map[string]string `json:"addons"`
	Files     []FileEntry       `json:"files"`
}

type PackageEntry struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Author      string                  `json:"author"`
	Repo        string                  `json:"repo"`
	Latest      string                  `json:"latest"`
	Versions    map[string]VersionEntry `json:"versions"`
}

func Fetch(name string) (*PackageEntry, error) {
	url := fmt.Sprintf("%s/packages/%s.json", registryBase, name)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry entry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	var entry PackageEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode registry entry: %w", err)
	}

	return &entry, nil
}

func Build(m *manifest.Manifest, fileURLs map[string]string) *PackageEntry {
	files := make([]FileEntry, 0, len(fileURLs))
	for name, url := range fileURLs {
		files = append(files, FileEntry{Name: name, URL: url})
	}

	versionEntry := VersionEntry{
		Skript:    m.Skript,
		Minecraft: m.Minecraft,
		Addons:    m.Addons,
		Files:     files,
	}

	return &PackageEntry{
		Name:        m.Name,
		Description: m.Description,
		Author:      m.Author,
		Repo:        m.Repo,
		Latest:      m.Version,
		Versions: map[string]VersionEntry{
			m.Version: versionEntry,
		},
	}
}

func Merge(existing *PackageEntry, m *manifest.Manifest, fileURLs map[string]string) *PackageEntry {
	files := make([]FileEntry, 0, len(fileURLs))
	for name, url := range fileURLs {
		files = append(files, FileEntry{Name: name, URL: url})
	}

	existing.Versions[m.Version] = VersionEntry{
		Skript:    m.Skript,
		Minecraft: m.Minecraft,
		Addons:    m.Addons,
		Files:     files,
	}
	existing.Latest = m.Version

	return existing
}

func Marshal(entry *PackageEntry) ([]byte, error) {
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package entry: %w", err)
	}

	return data, nil
}

func Encode(entry *PackageEntry, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entry)
}
