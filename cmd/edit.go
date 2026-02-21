package cmd

import (
	"errors"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var newName string
var newTarget string
var newTags []string
var newArgs []string
var newNote string

var EditBookmarkCmd = &cobra.Command{
	Use:   "edit <name> [flags]",
	Short: "edits the content of a bookmark",
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		pterm.Error.Println("usage: bm edit <name> [flags]")
		return errors.New("")
	}

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	name := args[0]
	found := false

	noteChanged := cmd.Flags().Changed("note")

	if newName == "" && newTarget == "" && len(newTags) == 0 && len(newArgs) == 0 && !noteChanged {
		pterm.Error.Println("no changes suppplied")
		return errors.New("no changes suppplied")
	}

	for i := range bookmarks {
		if bookmarks[i].Name == name {
			if newName != "" {
				bookmarks[i].Name = newName
			}
			if newTarget != "" {
				bookmarks[i].Target = newTarget
			}
			if len(newTags) != 0 {
				bookmarks[i].Tags = newTags
			}
			if len(newArgs) != 0 {
				bookmarks[i].Args = newArgs
				// Change type to app if args are being set
				if bookmarks[i].Type != "app" && bookmarks[i].Type != "url" {
					bookmarks[i].Type = "app"
				}
			}
			if noteChanged {
				bookmarks[i].Note = newNote
			}
			found = true
		}
	}

	if !found {
		pterm.Error.Printf("unable to find %s in bookmarks\n", name)
		return errors.New("unable to find bookmark")
	}

	err = internal.WriteBookmarks(&bookmarks)
	if err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("successfully edited %s\n", name)
	return nil
}

func init() {
	EditBookmarkCmd.Flags().StringVarP(&newName, "name", "n", "", "edits the name of the bookmark")
	EditBookmarkCmd.Flags().StringVarP(&newTarget, "target", "t", "", "edits the target of the bookmark")
	EditBookmarkCmd.Flags().StringSliceVar(&newTags, "tags", []string{}, "edits the tags of the bookmark")
	EditBookmarkCmd.Flags().StringSliceVarP(&newArgs, "args", "a", []string{}, "edits the arguments for app bookmarks")
	EditBookmarkCmd.Flags().StringVar(&newNote, "note", "", "sets a note on the bookmark")

	rootCmd.AddCommand(EditBookmarkCmd)
}
