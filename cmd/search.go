package cmd

import (
	"fmt"

	"github.com/skpm-dev/cli/internal/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search the skpm registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	results, err := api.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("No packages found matching %q\n", query)
		return nil
	}

	fmt.Printf("%-28s  %-12s  %-20s  %s\n", "NAME", "VERSION", "AUTHOR", "DESCRIPTION")
	fmt.Printf("%-28s  %-12s  %-20s  %s\n", "----", "-------", "------", "-----------")
	for _, p := range results {
		desc := p.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("%-28s  %-12s  %-20s  %s\n", p.Name, p.Latest, p.Author, desc)
	}
	return nil
}
