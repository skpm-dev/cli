package cmd

import (
	"fmt"
	"os"

	"github.com/skpm-dev/cli/internal/manifest"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your skpm.json without publishing",
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	m, err := manifest.Load("skpm.json")
	if err != nil {
		return err
	}

	if err := manifest.Validate(m); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	fmt.Println("skpm.json is valid")
	fmt.Println()
	fmt.Printf("  name:      %s\n", m.Name)
	fmt.Printf("  version:   %s\n", m.Version)
	fmt.Printf("  author:    %s\n", m.Author)
	fmt.Printf("  repo:      %s\n", m.Repo)
	if m.Skript != "" {
		fmt.Printf("  skript:    %s\n", m.Skript)
	}
	if m.Minecraft != "" {
		fmt.Printf("  minecraft: %s\n", m.Minecraft)
	}
	if len(m.Addons) > 0 {
		fmt.Printf("  addons:\n")
		for name, ver := range m.Addons {
			fmt.Printf("    %s %s\n", name, ver)
		}
	}
	if len(m.Dependencies) > 0 {
		fmt.Printf("  dependencies:\n")
		for name, ver := range m.Dependencies {
			fmt.Printf("    %s %s\n", name, ver)
		}
	}
	fmt.Printf("\nfiles (%d):\n", len(m.Files))
	for _, f := range m.Files {
		info, err := os.Stat(f)
		if err != nil {
			fmt.Printf("  %-40s  (not found on disk)\n", f)
		} else {
			fmt.Printf("  %-40s  %d bytes\n", f, info.Size())
		}
	}
	return nil
}
