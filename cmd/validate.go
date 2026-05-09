package cmd

import (
	"fmt"

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
	return nil
}
