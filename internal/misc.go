package internal

import (
	"fmt"

	"github.com/Shu-AFK/bm/internal/theme"
	"github.com/pterm/pterm"
)

func SelectCandidate(candidates []*Bookmark) *Bookmark {
	if len(candidates) == 0 {
		pterm.Warning.Println("No candidates available.")
		return nil
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	options := make([]string, len(candidates))
	for i, c := range candidates {
		options[i] = fmt.Sprintf("%s  (%s)", c.Name, c.Target)
	}

	printer := pterm.DefaultInteractiveSelect.
		WithMaxHeight(10).
		WithOptions(options).
		WithDefaultText(theme.GradientText("Select a bookmark"))

	printer.Selector = theme.GradientText("»")
	printer.SelectorStyle = theme.SecondaryStyle
	printer.TextStyle = theme.SecondaryStyle
	printer.OptionStyle = theme.PrimaryStyle

	choice, err := printer.Show()

	if err != nil {
		pterm.Error.Printf("Selection aborted: %v\n", err)
		return nil
	}

	for i, opt := range options {
		if opt == choice {
			return candidates[i]
		}
	}

	return nil
}
