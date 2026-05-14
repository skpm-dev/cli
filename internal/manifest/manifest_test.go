package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/skpm-dev/cli/internal/manifest"
)

func validManifest() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "economy",
		Description: "A test economy",
		Author:      "testuser",
		Version:     "1.0.0",
		Repo:        "https://github.com/test/economy",
		Files:       []string{"economy.sk"},
	}
}

func TestValidate_valid(t *testing.T) {
	if err := manifest.Validate(validManifest()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_missingFields(t *testing.T) {
	tests := []struct {
		field  string
		mutate func(*manifest.Manifest)
	}{
		{"name", func(m *manifest.Manifest) { m.Name = "" }},
		{"description", func(m *manifest.Manifest) { m.Description = "" }},
		{"author", func(m *manifest.Manifest) { m.Author = "" }},
		{"version", func(m *manifest.Manifest) { m.Version = "" }},
		{"files", func(m *manifest.Manifest) { m.Files = nil }},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			m := validManifest()
			tt.mutate(m)
			if err := manifest.Validate(m); err == nil {
				t.Fatalf("expected error for missing %s, got nil", tt.field)
			}
		})
	}
}

func TestSaveLoad_roundtrip(t *testing.T) {
	m := validManifest()
	m.Addons = map[string]string{"skript-reflect": ">=2.4"}

	tmp := filepath.Join(t.TempDir(), "skpm.json")

	if err := manifest.Save(m, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := manifest.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Name != m.Name {
		t.Errorf("name: got %q, want %q", loaded.Name, m.Name)
	}
	if loaded.Version != m.Version {
		t.Errorf("version: got %q, want %q", loaded.Version, m.Version)
	}
	if loaded.Addons["skript-reflect"] != m.Addons["skript-reflect"] {
		t.Errorf("addons: got %v, want %v", loaded.Addons, m.Addons)
	}
}

func TestLoad_missingFile(t *testing.T) {
	if _, err := manifest.Load("/nonexistent/skpm.json"); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
