package cmd

import (
	"strings"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

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

	if len(bookmarks) == 0 {
		pterm.Warning.Println("no bookmarks saved yet")
		return nil
	}

	data := pterm.TableData{
		{"Name", "Type", "Target", "Tags"},
	}

	for _, b := range bookmarks {
		data = append(data, []string{b.Name, b.Type, b.Target, strings.Join(b.Tags, ", ")})
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithHeaderRowSeparator("-").
		WithData(data).
		Render()

	return nil
}
func init() {
	rootCmd.AddCommand(ListBookmarksCmd)
}
