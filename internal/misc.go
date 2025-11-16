package internal

import (
	"fmt"

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

	choice, err := pterm.DefaultInteractiveSelect.
		WithMaxHeight(10).
		WithOptions(options).
		WithDefaultText("Select a bookmark").
		Show()

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
