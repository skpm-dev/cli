package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/skpm-dev/cli/internal/api"
	"github.com/skpm-dev/cli/internal/github"
	"github.com/skpm-dev/cli/internal/manifest"
	"github.com/skpm-dev/cli/internal/version"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a package to the skpm registry",
	RunE:  runPublish,
}

func init() {
	rootCmd.AddCommand(publishCmd)
}

func runPublish(cmd *cobra.Command, args []string) error {
	m, err := manifest.Load("skpm.json")
	if err != nil {
		return err
	}

	if err := manifest.Validate(m); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	token, err := github.Token()
	if err != nil {
		return err
	}

	existing, err := api.GetPackage(m.Name)
	if err != nil {
		return fmt.Errorf("could not check registry: %w", err)
	}

	if existing != nil {
		newVersion, err := promptVersionBump(existing.Latest)
		if err != nil {
			return err
		}
		m.Version = newVersion
	}

	files, err := readFiles(m.Files)
	if err != nil {
		return err
	}

	fmt.Printf("Publishing %s@%s...\n", m.Name, m.Version)

	if err := api.Publish(token, m, files); err != nil {
		if errors.Is(err, api.ErrVersionConflict) {
			newVersion, err := promptVersionBump(m.Version)
			if err != nil {
				return err
			}
			m.Version = newVersion
			fmt.Printf("Publishing %s@%s...\n", m.Name, m.Version)
			if err := api.Publish(token, m, files); err != nil {
				return fmt.Errorf("publish failed: %w", err)
			}
		} else {
			return fmt.Errorf("publish failed: %w", err)
		}
	}

	if err := manifest.Save(m, "skpm.json"); err != nil {
		fmt.Fprintf(os.Stderr, "warning: published successfully but could not update skpm.json: %v\n", err)
	}

	fmt.Printf("Published %s@%s\n", m.Name, m.Version)
	return nil
}

func readFiles(filenames []string) (map[string]string, error) {
	files := make(map[string]string, len(filenames))

	for _, name := range filenames {
		content, err := os.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("could not read file %s: %w", name, err)
		}
		files[name] = string(content)
	}

	return files, nil
}

func promptVersionBump(current string) (string, error) {
	v, err := version.Parse(current)
	if err != nil {
		return "", fmt.Errorf("could not parse current version %q: %w", current, err)
	}

	patch := version.Bump(v, version.BumpPatch)
	minor := version.Bump(v, version.BumpMinor)
	major := version.Bump(v, version.BumpMajor)

	fmt.Printf("\nFound existing package at version %s\n", current)
	fmt.Printf("What type of release is this?\n")
	fmt.Printf("  [1] patch — bug fixes        (%s → %s)\n", current, patch)
	fmt.Printf("  [2] minor — new features     (%s → %s)\n", current, minor)
	fmt.Printf("  [3] major — breaking changes (%s → %s)\n", current, major)
	fmt.Print("\nEnter choice: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("could not read input: %w", err)
	}

	switch strings.TrimSpace(input) {
	case "1":
		return patch.String(), nil
	case "2":
		return minor.String(), nil
	case "3":
		return major.String(), nil
	default:
		return "", fmt.Errorf("invalid choice %q, expected 1, 2 or 3", strings.TrimSpace(input))
	}
}
