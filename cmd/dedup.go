package cmd

import (
	"fmt"

	"github.com/Shu-AFK/bm/internal"
	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var DedupCmd = &cobra.Command{
	Use:   "dedup",
	Short: "Find and interactively remove duplicate bookmark targets",
	RunE:  dedupCmd,
}

func dedupCmd(cmd *cobra.Command, args []string) error {
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	targetMap := make(map[string][]int)
	for i, b := range bookmarks {
		targetMap[b.Target] = append(targetMap[b.Target], i)
	}

	type dupGroup struct {
		target  string
		indices []int
	}
	var dups []dupGroup
	for target, indices := range targetMap {
		if len(indices) > 1 {
			dups = append(dups, dupGroup{target, indices})
		}
	}

	if len(dups) == 0 {
		pterm.Success.Println("no duplicate targets found")
		return nil
	}

	pterm.Printf("found %d target(s) with duplicates\n\n", len(dups))

	toDelete := make(map[int]bool)

	for _, dup := range dups {
		pterm.Printf("Duplicate target: %s\n", theme.GradientText(dup.target))

		options := make([]string, len(dup.indices))
		for i, idx := range dup.indices {
			b := bookmarks[idx]
			options[i] = fmt.Sprintf("%s  (created: %s)", b.Name, b.CreatedAt)
		}

		printer := pterm.DefaultInteractiveSelect.
			WithMaxHeight(10).
			WithOptions(options).
			WithDefaultText(theme.GradientText("Select bookmark to KEEP (others will be deleted)"))

		printer.Selector = theme.SecondaryStyle.Sprint("»")
		printer.SelectorStyle = theme.SecondaryStyle
		printer.TextStyle = theme.SecondaryStyle
		printer.OptionStyle = theme.SecondaryStyle

		choice, err := printer.Show()
		if err != nil {
			pterm.Warning.Printf("skipping duplicate group: %v\n", err)
			continue
		}

		for i, opt := range options {
			if opt != choice {
				toDelete[dup.indices[i]] = true
			}
		}
		pterm.Println()
	}

	if len(toDelete) == 0 {
		pterm.Info.Println("no bookmarks selected for deletion")
		return nil
	}

	confirmText := fmt.Sprintf("Delete %d duplicate bookmark(s)?", len(toDelete))
	confirmed, err := pterm.DefaultInteractiveConfirm.
		WithDefaultText(confirmText).
		WithDefaultValue(false).
		Show()
	if err != nil || !confirmed {
		pterm.Info.Println("operation cancelled")
		return nil
	}

	var kept []internal.Bookmark
	for i, b := range bookmarks {
		if !toDelete[i] {
			kept = append(kept, b)
		}
	}
	if kept == nil {
		kept = []internal.Bookmark{}
	}

	if err := internal.WriteBookmarks(&kept); err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("removed %d duplicate bookmark(s)\n", len(toDelete))
	return nil
}

func init() {
	rootCmd.AddCommand(DedupCmd)
}
