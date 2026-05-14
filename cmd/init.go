package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/skpm-dev/cli/internal/github"
	"github.com/spf13/cobra"
)

// reNonPackageChars matches any character not allowed in a registry package
// name. Used to coerce a directory name into a publishable slug.
var reNonPackageChars = regexp.MustCompile(`[^a-z0-9-]+`)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a skpm.json in the current directory",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	const path = "skpm.json"

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	// Default name from the current directory name, coerced to the registry's
	// allowed character set (^[a-z][a-z0-9-]{1,38}$).
	name := "my-package"
	if cwd, err := os.Getwd(); err == nil {
		slug := reNonPackageChars.ReplaceAllString(strings.ToLower(filepath.Base(cwd)), "-")
		slug = strings.Trim(slug, "-")
		if len(slug) > 39 {
			slug = slug[:39]
		}
		if slug != "" && slug[0] >= 'a' && slug[0] <= 'z' {
			name = slug
		}
	}

	// Try to pre-fill author and repo from the GitHub token.
	author := "your-github-username"
	repo := author + "/" + name
	if tok, err := github.Token(); err == nil {
		if user, err := github.GetAuthenticatedUser(tok); err == nil {
			author = user.Login
			repo = user.Login + "/" + name
		}
	}

	m := map[string]any{
		"name":         name,
		"description":  "A short description of your package",
		"author":       author,
		"version":      "1.0.0",
		"repo":         repo,
		"skript":       ">=2.7",
		"minecraft":    ">=1.20",
		"addons":       map[string]string{},
		"dependencies": map[string]string{},
		"files":        []string{"main.sk"},
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode skpm.json: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}

	fmt.Printf("Created skpm.json for %s\n", name)
	if author == "your-github-username" {
		fmt.Println("Tip: set SKPM_GITHUB_TOKEN to auto-fill author and repo next time.")
	}
	return nil
}
