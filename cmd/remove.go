package cmd

import (
	"fmt"
	"os"

	"github.com/skpm-dev/cli/internal/api"
	"github.com/spf13/cobra"
)

var removeReason string

var removeCmd = &cobra.Command{
	Use:   "remove <package> [version]",
	Short: "Remove or yank a package from the registry (admin only)",
	Long: `Remove a package or yank a single version from the registry.

Requires the SKPM_ADMIN_TOKEN environment variable to be set.

  skpm remove economy              # hard-remove entire package
  skpm remove economy 1.0.0        # yank a single version
  skpm remove economy --reason "contains malicious code"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runRemove,
}

func init() {
	removeCmd.Flags().StringVar(&removeReason, "reason", "", "Reason for removal (recorded in the registry)")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	token := os.Getenv("SKPM_ADMIN_TOKEN")
	if token == "" {
		return fmt.Errorf("SKPM_ADMIN_TOKEN is not set")
	}

	name := args[0]

	if len(args) == 2 {
		version := args[1]
		fmt.Printf("Yanking %s@%s...\n", name, version)
		if err := api.Yank(token, name, version, removeReason); err != nil {
			return fmt.Errorf("yank failed: %w", err)
		}
		fmt.Printf("Yanked %s@%s\n", name, version)
		return nil
	}

	fmt.Printf("Removing %s...\n", name)
	if err := api.Remove(token, name, removeReason); err != nil {
		return fmt.Errorf("remove failed: %w", err)
	}
	fmt.Printf("Removed %s\n", name)
	return nil
}
