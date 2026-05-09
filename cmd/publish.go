package cmd

import (
	"fmt"

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
	fmt.Println("publish: not yet implemented")
	return nil
}
