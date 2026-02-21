package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var CopyBookmarkCmd = &cobra.Command{
	Use:   "copy <name>",
	Short: "Copy the bookmark target to the clipboard",
	RunE:  copyCmd,
}

func copyCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm copy <name>")
		return errors.New("wrong usage")
	}

	name := args[0]
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	match := internal.GetClosestMatch(bookmarks, name)
	if match.Kind == internal.MatchNone {
		pterm.Error.Printf("could not find match for %s\n", name)
		return errors.New("no match found")
	}

	var bm *internal.Bookmark
	if match.BM != nil {
		bm = match.BM
	} else {
		bm = internal.SelectCandidate(match.Candidates)
		if bm == nil {
			return nil
		}
	}

	if err := copyToClipboard(bm.Target); err != nil {
		pterm.Error.Printf("failed to copy to clipboard: %v\n", err)
		return err
	}

	pterm.Success.Printf("copied %s to clipboard\n", bm.Target)
	return nil
}

func copyToClipboard(text string) error {
	type clipCmd struct {
		name string
		args []string
	}

	candidates := []clipCmd{
		{"wl-copy", []string{}},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
		{"pbcopy", []string{}},
		{"clip", []string{}},
	}

	for _, c := range candidates {
		if _, err := exec.LookPath(c.name); err != nil {
			continue
		}
		cmd := exec.Command(c.name, c.args...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}

	return errors.New("no clipboard utility found (install wl-copy, xclip, xsel, pbcopy, or clip)")
}

func init() {
	rootCmd.AddCommand(CopyBookmarkCmd)
}
