package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skpm-dev/cli/internal/github"
	"github.com/skpm-dev/cli/internal/manifest"
	"github.com/skpm-dev/cli/internal/registry"
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

	existing, err := registry.Fetch(m.Name)
	if err != nil {
		return fmt.Errorf("could not check registry: %w", err)
	}

	if existing != nil {
		newVersion, err := promptVersionBump(existing.Latest)
		if err != nil {
			return err
		}
		m.Version = newVersion
		if err := manifest.Save(m, "skpm.json"); err != nil {
			return fmt.Errorf("could not update skpm.json: %w", err)
		}
		fmt.Printf("Bumped version to %s in skpm.json\n", m.Version)
	}

	client, err := github.NewClient()
	if err != nil {
		return err
	}

	parts := strings.SplitN(m.Repo, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("repo must be in format owner/repo, got %q", m.Repo)
	}
	owner, repo := parts[0], parts[1]

	tag := "v" + m.Version
	fmt.Printf("Creating GitHub release %s...\n", tag)

	release, err := client.CreateRelease(owner, repo, tag, m.Version)
	if err != nil {
		return fmt.Errorf("could not create release: %w", err)
	}

	fileURLs := make(map[string]string)
	for _, file := range m.Files {
		fmt.Printf("Uploading %s...\n", file)
		url, err := client.UploadReleaseAsset(release, file)
		if err != nil {
			return fmt.Errorf("could not upload %s: %w", file, err)
		}
		fileURLs[file] = url
	}

	var entry *registry.PackageEntry
	if existing != nil {
		entry = registry.Merge(existing, m, fileURLs)
	} else {
		entry = registry.Build(m, fileURLs)
	}

	entryJSON, err := registry.Marshal(entry)
	if err != nil {
		return err
	}

	fmt.Println("Opening PR to skpm-dev/registry...")
	fmt.Printf("\nPackage entry that will be submitted:\n\n%s\n", string(entryJSON))
	fmt.Println("(PR creation coming in next step)")

	return nil
}

func promptVersionBump(current string) (string, error) {
	v, err := version.Parse(current)
	if err != nil {
		return "", fmt.Errorf("could not parse current version %q: %w", current, err)
	}

	patch := v.Bump(version.BumpPatch)
	minor := v.Bump(version.BumpMinor)
	major := v.Bump(version.BumpMajor)

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
