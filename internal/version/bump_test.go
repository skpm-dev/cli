package version_test

import (
	"testing"

	"github.com/skpm-dev/cli/internal/version"
)

func TestParse_valid(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"0.0.1", "0.0.1"},
		{"1.0.0-beta.1", "1.0.0-beta.1"},
		{"2.0.0", "2.0.0"},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			v, err := version.Parse(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.String() != c.expected {
				t.Fatalf("got %s, want %s", v.String(), c.expected)
			}
		})
	}
}

func TestParse_invalid(t *testing.T) {
	cases := []string{"not-a-version", "", "abc", "1.2.3.4.5"}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			if _, err := version.Parse(c); err == nil {
				t.Fatalf("expected error for %q, got nil", c)
			}
		})
	}
}

func TestBump(t *testing.T) {
	v, _ := version.Parse("1.2.3")

	cases := []struct {
		bumpType version.BumpType
		expected string
	}{
		{version.BumpPatch, "1.2.4"},
		{version.BumpMinor, "1.3.0"},
		{version.BumpMajor, "2.0.0"},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			result := version.Bump(v, c.bumpType)
			if result.String() != c.expected {
				t.Fatalf("got %s, want %s", result.String(), c.expected)
			}
		})
	}
}

func TestBump_doesNotMutateOriginal(t *testing.T) {
	v, _ := version.Parse("1.2.3")
	version.Bump(v, version.BumpMajor)
	if v.String() != "1.2.3" {
		t.Fatalf("original version was mutated: got %s", v.String())
	}
}
