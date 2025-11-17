package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var ImportBookmarksCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Import bookmarks from a file, replacing the current set",
	RunE:  importBookmarks,
}

func importBookmarks(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("bm import <path>")
		return errors.New("wrong usage")
	}

	src := args[0]

	data, err := os.ReadFile(src)
	if err != nil {
		pterm.Error.Printf("unable to read %s: %v\n", src, err)
		return err
	}

	var bookmarks []internal.Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		pterm.Error.Printf("invalid bookmark file: %v\n", err)
		return err
	}

	if err := internal.WriteBookmarks(&bookmarks); err != nil {
		pterm.Error.Printf("unable to save imported bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("imported %d bookmarks from %s\n", len(bookmarks), src)
	return nil
}

func init() {
	rootCmd.AddCommand(ImportBookmarksCmd)
}
