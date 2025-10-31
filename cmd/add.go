package cmd

import "github.com/spf13/cobra"

var tags []string

var AddBookmarkCmd = &cobra.Command{
	Use:   "add <name> <path / url>",
	Short: "Add a new bookmark",

	Args: cobra.ExactArgs(2),
	RunE: exec,
}

func exec(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	AddBookmarkCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "All tags associated with that bookmark, comma seperated")

	rootCmd.AddCommand(AddBookmarkCmd)
}
