package cmd

import (
	"fmt"
	"sort"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show usage statistics",
	RunE:  statsCmd,
}

func statsCmd(cmd *cobra.Command, args []string) error {
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	if len(bookmarks) == 0 {
		pterm.Warning.Println("no bookmarks saved yet")
		return nil
	}

	typeCounts := make(map[string]int)
	pinnedCount := 0
	tagCounts := make(map[string]int)

	for _, b := range bookmarks {
		typeCounts[b.Type]++
		if b.Pinned {
			pinnedCount++
		}
		for _, t := range b.Tags {
			tagCounts[t]++
		}
	}

	sorted := make([]internal.Bookmark, len(bookmarks))
	copy(sorted, bookmarks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].AccessCount > sorted[j].AccessCount
	})

	type tagEntry struct {
		name  string
		count int
	}
	var tagList []tagEntry
	for name, count := range tagCounts {
		tagList = append(tagList, tagEntry{name, count})
	}
	sort.Slice(tagList, func(i, j int) bool {
		if tagList[i].count != tagList[j].count {
			return tagList[i].count > tagList[j].count
		}
		return tagList[i].name < tagList[j].name
	})

	pterm.Println(theme.GradientText("OVERVIEW"))
	pterm.Println(theme.GradientText("────────"))
	pterm.Printf("  Total bookmarks : %d\n", len(bookmarks))
	pterm.Printf("  Pinned          : %d\n", pinnedCount)
	pterm.Println()

	pterm.Println(theme.GradientText("BY TYPE"))
	pterm.Println(theme.GradientText("───────"))
	typeData := pterm.TableData{{
		theme.GradientText("Type"),
		theme.GradientText("Count"),
	}}
	for _, t := range []string{"url", "folder", "file", "app"} {
		if c, ok := typeCounts[t]; ok {
			typeData = append(typeData, []string{t, fmt.Sprintf("%d", c)})
		}
	}
	for t, c := range typeCounts {
		if t != "url" && t != "folder" && t != "file" && t != "app" {
			typeData = append(typeData, []string{t, fmt.Sprintf("%d", c)})
		}
	}
	pterm.DefaultTable.
		WithHasHeader().
		WithSeparator("│").
		WithHeaderRowSeparator("─").
		WithBoxed(true).
		WithData(typeData).
		Render()
	pterm.Println()

	pterm.Println(theme.GradientText("TOP 5 MOST ACCESSED"))
	pterm.Println(theme.GradientText("───────────────────"))
	topData := pterm.TableData{{
		theme.GradientText("Name"),
		theme.GradientText("Target"),
		theme.GradientText("Access Count"),
	}}
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		b := sorted[i]
		if b.AccessCount == 0 {
			break
		}
		topData = append(topData, []string{
			b.Name,
			b.Target,
			fmt.Sprintf("%d", b.AccessCount),
		})
	}
	if len(topData) == 1 {
		pterm.Println("  No bookmarks accessed yet.")
	} else {
		pterm.DefaultTable.
			WithHasHeader().
			WithSeparator("│").
			WithHeaderRowSeparator("─").
			WithBoxed(true).
			WithData(topData).
			Render()
	}
	pterm.Println()

	if len(tagList) > 0 {
		pterm.Println(theme.GradientText("TOP 10 TAGS"))
		pterm.Println(theme.GradientText("───────────"))
		tagData := pterm.TableData{{
			theme.GradientText("Tag"),
			theme.GradientText("Count"),
		}}
		tagLimit := 10
		if len(tagList) < tagLimit {
			tagLimit = len(tagList)
		}
		for i := 0; i < tagLimit; i++ {
			tagData = append(tagData, []string{tagList[i].name, fmt.Sprintf("%d", tagList[i].count)})
		}
		pterm.DefaultTable.
			WithHasHeader().
			WithSeparator("│").
			WithHeaderRowSeparator("─").
			WithBoxed(true).
			WithData(tagData).
			Render()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(StatsCmd)
}
