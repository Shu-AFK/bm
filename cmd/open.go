package cmd

import (
	"errors"
	"os/exec"
	"runtime"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var OpenBookmarkCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "opens the bookmark with the default browser/editor or cd's into the directory",
	RunE:  openCmd,
}

func openCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("bm open <name>")
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

	switch match.Kind {
	case internal.MatchExact, internal.MatchPrefix, internal.MatchFuzzy:
		var chosen *internal.Bookmark
		if match.BM != nil {
			chosen = match.BM
		} else {
			chosen = internal.SelectCandidate(match.Candidates)
			if chosen == nil {
				return nil
			}
		}

		if err := open(chosen, name); err != nil {
			return err
		}

		// Update access tracking
		chosen.LastAccessed = time.Now().Format(time.RFC3339)
		chosen.AccessCount++
		if writeErr := internal.WriteBookmarks(&bookmarks); writeErr != nil {
			pterm.Warning.Printf("could not update access stats: %v\n", writeErr)
		}
		return nil

	case internal.MatchNone:
		pterm.Error.Printf("could not find bookmark for %s\n", name)
		return errors.New("no match found")
	}

	return errors.New("unreachable")
}

func open(bm *internal.Bookmark, name string) error {
	target := bm.Target
	var err error

	if bm.Type == "folder" {
		pterm.Printf("cd %s\n", target)
		return nil
	}

	if bm.Type == "app" {
		err = exec.Command(target, bm.Args...).Start()
		if err != nil {
			pterm.Error.Printf("unable to run %s: %v\n", name, err)
			return err
		}
		pterm.Printf("launched %s\n", target)
		return nil
	}

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
	case "darwin":
		err = exec.Command("open", target).Start()
	default:
		err = exec.Command("xdg-open", target).Start()
	}

	if err != nil {
		pterm.Error.Printf("unable to open %s: %v", name, err)
		return err
	}

	pterm.Printf("opened %s\n", target)

	return nil
}

func init() {
	rootCmd.AddCommand(OpenBookmarkCmd)
}
