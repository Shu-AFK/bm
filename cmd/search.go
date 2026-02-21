package cmd

import (
	"errors"
	"sort"
	"strings"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var SearchBookmarksCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Full-text search across name, tags, and notes",
	RunE:  searchCmd,
}

type searchResult struct {
	bm    internal.Bookmark
	score int
}

func searchCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		pterm.Error.Println("usage: bm search <query>")
		return errors.New("wrong usage")
	}

	query := strings.ToLower(strings.Join(args, " "))

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	var results []searchResult
	for _, b := range bookmarks {
		score := scoreBookmark(b, query)
		if score > 0 {
			results = append(results, searchResult{bm: b, score: score})
		}
	}

	if len(results) == 0 {
		pterm.Warning.Printf("no bookmarks found matching %q\n", query)
		return nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	headers := []string{
		theme.GradientText("Name"),
		theme.GradientText("Type"),
		theme.GradientText("Target"),
		theme.GradientText("Tags"),
		theme.GradientText("Note"),
	}

	data := pterm.TableData{headers}
	for _, r := range results {
		data = append(data, []string{
			r.bm.Name,
			r.bm.Type,
			r.bm.Target,
			strings.Join(r.bm.Tags, ", "),
			truncate(r.bm.Note, 40),
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

func scoreBookmark(b internal.Bookmark, query string) int {
	score := 0
	nameLower := strings.ToLower(b.Name)

	if nameLower == query {
		score += 10
	} else if strings.Contains(nameLower, query) {
		score += 5
	}

	for _, t := range b.Tags {
		if strings.Contains(strings.ToLower(t), query) {
			score += 3
			break
		}
	}

	if b.Note != "" && strings.Contains(strings.ToLower(b.Note), query) {
		score += 1
	}

	return score
}

func init() {
	rootCmd.AddCommand(SearchBookmarksCmd)
}
