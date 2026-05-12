package cmd

import (
	"fmt"
	"sort"

	"github.com/skpm-dev/cli/internal/api"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <package>",
	Short: "Show details about a package in the skpm registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	name := args[0]

	pkg, err := api.GetPackage(name)
	if err != nil {
		return fmt.Errorf("could not fetch package: %w", err)
	}
	if pkg == nil {
		return fmt.Errorf("package %q not found", name)
	}

	fmt.Printf("  name:        %s\n", pkg.Name)
	fmt.Printf("  description: %s\n", pkg.Description)
	fmt.Printf("  author:      %s\n", pkg.Author)
	fmt.Printf("  latest:      %s\n", pkg.Latest)

	if v, ok := pkg.Versions[pkg.Latest]; ok {
		if v.Skript != "" {
			fmt.Printf("  skript:      %s\n", v.Skript)
		}
		if v.Minecraft != "" {
			fmt.Printf("  minecraft:   %s\n", v.Minecraft)
		}
		if len(v.Addons) > 0 {
			fmt.Printf("  addons:\n")
			for addon, ver := range v.Addons {
				fmt.Printf("    %s %s\n", addon, ver)
			}
		}
		if len(v.Files) > 0 {
			fmt.Printf("  files:\n")
			for _, f := range v.Files {
				fmt.Printf("    %s\n", f.Name)
			}
		}
	}

	// Version history
	versions := make([]string, 0, len(pkg.Versions))
	for v := range pkg.Versions {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	fmt.Printf("\nversions (%d):\n", len(versions))
	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		entry := pkg.Versions[v]
		if entry.Yanked {
			fmt.Printf("  %s  [yanked]\n", v)
		} else if v == pkg.Latest {
			fmt.Printf("  %s  (latest)\n", v)
		} else {
			fmt.Printf("  %s\n", v)
		}
	}
	return nil
}
