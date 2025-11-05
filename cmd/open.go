package cmd

import "github.com/spf13/cobra"

var OpenBookmarkCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "opens the bookmark with the default browser/editor or cd's into the directory",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(OpenBookmarkCmd)
}
