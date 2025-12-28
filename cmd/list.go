package cmd

import (
	"errors"
	"strings"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var tag string
var sortmode string

var ListBookmarksCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all bookmarks",

	RunE: list,
}

func list(cmd *cobra.Command, args []string) error {
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	if tag != "" {
		bookmarks = internal.SearchByTag(bookmarks, tag)
	}

	if len(bookmarks) == 0 && tag == "" {
		pterm.Warning.Println("no bookmarks saved yet")
		return nil
	} else if len(bookmarks) == 0 {
		pterm.Warning.Println("no bookmarks with this tag found")
		return nil
	}

	if sortmode != "" {
		var mode internal.SortMode

		switch sortmode {
		case "name":
			mode = internal.SortByName
		case "time":
			mode = internal.SortByCreated
		case "target":
			mode = internal.SortByTarget
		case "type":
			mode = internal.SortByType
		default:
			pterm.Error.Printf("%s is not a valid sort mode, please choose from the following: name, time, target, type.\n", sortmode)
			return errors.New("wrong sort mode")
		}

		internal.SortBookmarks(bookmarks, mode)
	}

	headers := []string{
		theme.GradientText("Name"),
		theme.GradientText("Type"),
		theme.GradientText("Target"),
		theme.GradientText("Args"),
		theme.GradientText("Tags"),
	}

	data := pterm.TableData{headers}

	for _, b := range bookmarks {
		data = append(data, []string{
			b.Name,
			b.Type,
			b.Target,
			strings.Join(b.Args, " "),
			strings.Join(b.Tags, ", "),
		})
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithSeparator("│").
		WithHeaderRowSeparator("─").
		WithBoxed(true).
		WithData(data).
		Render()

	return nil
}
func init() {
	ListBookmarksCmd.Flags().StringVarP(&tag, "tag", "t", "", "filter by tag")
	ListBookmarksCmd.Flags().StringVarP(&sortmode, "sort", "s", "", "sort by name, time, target or type")

	rootCmd.AddCommand(ListBookmarksCmd)
}
