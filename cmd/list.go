package cmd

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var tag string
var sortmode string
var pinnedOnly bool

var ListBookmarksCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all bookmarks",

	RunE: list,
}

func truncate(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max-1]) + "…"
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

	if pinnedOnly {
		var pinned []internal.Bookmark
		for _, b := range bookmarks {
			if b.Pinned {
				pinned = append(pinned, b)
			}
		}
		bookmarks = pinned
	}

	if len(bookmarks) == 0 && tag == "" && !pinnedOnly {
		pterm.Warning.Println("no bookmarks saved yet")
		return nil
	} else if len(bookmarks) == 0 {
		pterm.Warning.Println("no bookmarks with this filter found")
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

	// Pinned bookmarks always first
	var pinned, unpinned []internal.Bookmark
	for _, b := range bookmarks {
		if b.Pinned {
			pinned = append(pinned, b)
		} else {
			unpinned = append(unpinned, b)
		}
	}
	bookmarks = append(pinned, unpinned...)

	headers := []string{
		theme.GradientText("Name"),
		theme.GradientText("Type"),
		theme.GradientText("Target"),
		theme.GradientText("Args"),
		theme.GradientText("Tags"),
		theme.GradientText("Note"),
	}

	data := pterm.TableData{headers}

	for _, b := range bookmarks {
		name := b.Name
		if b.Pinned {
			name = "★ " + name
		}
		data = append(data, []string{
			name,
			b.Type,
			b.Target,
			strings.Join(b.Args, " "),
			strings.Join(b.Tags, ", "),
			truncate(b.Note, 40),
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
	ListBookmarksCmd.Flags().BoolVarP(&pinnedOnly, "pinned", "p", false, "show only pinned bookmarks")

	rootCmd.AddCommand(ListBookmarksCmd)
}
