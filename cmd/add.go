package cmd

import (
	"errors"
	"net/url"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var tags []string

var AddBookmarkCmd = &cobra.Command{
	Use:   "add <name> <path/url>",
	Short: "Add a new bookmark",

	RunE: exec,
}

func exec(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		pterm.Error.Println("usage: bm add <name> <path/url> [Flags]")
		return nil
	}

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	name := args[0]
	target := args[1]

	for _, bookmark := range bookmarks {
		if name == bookmark.Name {
			pterm.Error.Println("a bookmark with this name already exists")
			return errors.New("name already exists")
		}
	}

	bookmarks = append(bookmarks, internal.Bookmark{
		Name:      name,
		Target:    target,
		Type:      getBookmarkType(target),
		Tags:      tags,
		CreatedAt: time.Now().String(),
	})

	err = internal.WriteBookmarks(&bookmarks)
	if err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("successfully added bookmark %s\n", name)

	return nil
}

func getBookmarkType(target string) string {
	if u, err := url.Parse(target); err == nil && u.Scheme != "" && u.Host != "" {
		return "url"
	}

	return "path"
}

func init() {
	AddBookmarkCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "All tags associated with that bookmark, comma seperated")

	rootCmd.AddCommand(AddBookmarkCmd)
}
