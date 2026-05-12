package models

type PackageSummary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Latest      string `json:"latest"`
	Downloads   int64  `json:"downloads"`
}

type FileEntry struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

type VersionEntry struct {
	Skript    string            `json:"skript"`
	Minecraft string            `json:"minecraft"`
	Addons    map[string]string `json:"addons"`
	Files     []FileEntry       `json:"files"`
	Yanked    bool              `json:"yanked,omitempty"`
	Downloads int64             `json:"downloads,omitempty"`
}

type Package struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Author      string                  `json:"author"`
	Latest      string                  `json:"latest"`
	Downloads   int64                   `json:"downloads,omitempty"`
	Versions    map[string]VersionEntry `json:"versions"`
}
