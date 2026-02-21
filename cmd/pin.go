package cmd

import (
	"errors"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var PinBookmarkCmd = &cobra.Command{
	Use:   "pin <name>",
	Short: "Toggle the pinned state of a bookmark",
	RunE:  pinCmd,
}

func pinCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm pin <name>")
		return errors.New("wrong usage")
	}

	name := args[0]
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	found := false
	for i := range bookmarks {
		if bookmarks[i].Name == name {
			bookmarks[i].Pinned = !bookmarks[i].Pinned
			if bookmarks[i].Pinned {
				pterm.Success.Printf("pinned %s\n", name)
			} else {
				pterm.Success.Printf("unpinned %s\n", name)
			}
			found = true
			break
		}
	}

	if !found {
		pterm.Error.Printf("bookmark %s not found\n", name)
		return errors.New("bookmark not found")
	}

	if err := internal.WriteBookmarks(&bookmarks); err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(PinBookmarkCmd)
}
