package cmd

import (
	"errors"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var RemoveBookmarkCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a bookmark",

	RunE: remove,
}

func remove(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm remove <name>")
		return errors.New("wrong usage")
	}

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read booksmarks: %v\n", err)
		return err
	}

	success := false
	for i, bookmark := range bookmarks {
		if bookmark.Name == args[0] {
			bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
			success = true
			break
		}
	}

	if success == false {
		pterm.Error.Printf("name %s not found\n", args[0])
		return errors.New("name not found")
	}

	err = internal.WriteBookmarks(&bookmarks)
	if err != nil {
		pterm.Error.Printf("Unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("Successfully removed bookemark: %s", args[0])
	return nil
}

func init() {
	rootCmd.AddCommand(RemoveBookmarkCmd)
}
