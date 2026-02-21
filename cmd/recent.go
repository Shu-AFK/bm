package cmd

import (
	"sort"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var recentCount int

var RecentBookmarksCmd = &cobra.Command{
	Use:   "recent",
	Short: "Show the most recently accessed bookmarks",
	RunE:  recentCmd,
}

func recentCmd(cmd *cobra.Command, args []string) error {
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	var accessed []internal.Bookmark
	for _, b := range bookmarks {
		if b.LastAccessed != "" {
			accessed = append(accessed, b)
		}
	}

	if len(accessed) == 0 {
		pterm.Warning.Println("no recently accessed bookmarks")
		return nil
	}

	sort.Slice(accessed, func(i, j int) bool {
		ti, errI := time.Parse(time.RFC3339, accessed[i].LastAccessed)
		tj, errJ := time.Parse(time.RFC3339, accessed[j].LastAccessed)
		if errI != nil || errJ != nil {
			return accessed[i].LastAccessed > accessed[j].LastAccessed
		}
		return ti.After(tj)
	})

	if recentCount > 0 && len(accessed) > recentCount {
		accessed = accessed[:recentCount]
	}

	headers := []string{
		theme.GradientText("Name"),
		theme.GradientText("Type"),
		theme.GradientText("Target"),
		theme.GradientText("Last Accessed"),
	}

	data := pterm.TableData{headers}
	for _, b := range accessed {
		data = append(data, []string{
			b.Name,
			b.Type,
			b.Target,
			b.LastAccessed,
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
	RecentBookmarksCmd.Flags().IntVarP(&recentCount, "count", "c", 10, "number of recent bookmarks to show")
	rootCmd.AddCommand(RecentBookmarksCmd)
}
