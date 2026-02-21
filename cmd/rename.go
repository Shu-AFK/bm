package cmd

import (
	"errors"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var RenameBookmarkCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a bookmark",
	RunE:  renameCmd,
}

func renameCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		pterm.Error.Println("usage: bm rename <old> <new>")
		return errors.New("wrong usage")
	}

	oldName := args[0]
	newNameArg := args[1]

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	for _, b := range bookmarks {
		if b.Name == newNameArg {
			pterm.Error.Printf("a bookmark named %s already exists\n", newNameArg)
			return errors.New("name already exists")
		}
	}

	found := false
	for i := range bookmarks {
		if bookmarks[i].Name == oldName {
			bookmarks[i].Name = newNameArg
			found = true
			break
		}
	}

	if !found {
		pterm.Error.Printf("bookmark %s not found\n", oldName)
		return errors.New("bookmark not found")
	}

	if err := internal.WriteBookmarks(&bookmarks); err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	groups, err := internal.ReadGroups()
	if err == nil {
		updated := false
		for i := range groups {
			for j, name := range groups[i].Bookmarks {
				if name == oldName {
					groups[i].Bookmarks[j] = newNameArg
					updated = true
				}
			}
		}
		if updated {
			if writeErr := internal.WriteGroups(&groups); writeErr != nil {
				pterm.Warning.Printf("bookmark renamed but could not update group references: %v\n", writeErr)
			}
		}
	}

	pterm.Success.Printf("renamed %s to %s\n", oldName, newNameArg)
	return nil
}

func init() {
	rootCmd.AddCommand(RenameBookmarkCmd)
}
