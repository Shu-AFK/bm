package cmd

import (
	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var forceFlag bool

var ClearBookmarksCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all bookmarks",
	RunE:  clear,
}

func clear(cmd *cobra.Command, args []string) error {
	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	if len(bookmarks) == 0 {
		pterm.Info.Println("no bookmarks to clear")
		return nil
	}

	if !forceFlag {
		result, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Are you sure you want to remove all bookmarks?").
			Show()
		if !result {
			pterm.Info.Println("aborted")
			return nil
		}
	}

	bookmarks = []internal.Bookmark{}
	err = internal.WriteBookmarks(&bookmarks)
	if err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Println("all bookmarks cleared")
	return nil
}

func init() {
	ClearBookmarksCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(ClearBookmarksCmd)
}
