package cmd

import (
	"errors"
	"os/exec"
	"runtime"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var OpenBookmarkCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "opens the bookmark with the default browser/editor or cd's into the directory",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
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

	var found *internal.Bookmark
	for _, bookmark := range bookmarks {
		if bookmark.Name == name {
			found = &bookmark
			break
		}
	}

	if found == nil {
		pterm.Error.Printf("unable to find bookmark %s\n", name)
		return errors.New("not found")
	}

	if found.Type == "path" {
		pterm.Error.Println("bookmark is path, currently not implemented")
		return errors.New("unimplemented")
	}

	target := found.Target

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

	return nil
}

func init() {
	rootCmd.AddCommand(OpenBookmarkCmd)
}
