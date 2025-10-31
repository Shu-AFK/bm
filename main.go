package main

import (
	"os"

	"github.com/Shu-AFK/bm/cmd"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bm",
	Short: "CLI Bookmark Manager: save and open links fast",
}

func main() {
	rootCmd.AddCommand(cmd.AddBookmarkCmd)

	err := rootCmd.Execute()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
