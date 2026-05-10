package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a template skpm.json in the current directory",
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

	template := `{
  "name": "my-package",
  "description": "A short description of your package",
  "author": "your-github-username",
  "version": "1.0.0",
  "repo": "your-github-username/your-repo",
  "skript": ">=2.7",
  "minecraft": ">=1.20",
  "addons": {},
  "files": ["main.sk"]
}
`

	if err := os.WriteFile(path, []byte(template), 0644); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}

	fmt.Println("Created skpm.json")
	return nil
}
