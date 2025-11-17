package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var ExportBookmarksCmd = &cobra.Command{
	Use:   "export <path>",
	Short: "Export saved bookmarks to a file",
	RunE:  export,
}

func export(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("bm export <path>")
		return errors.New("wrong usage")
	}

	dest := args[0]

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		pterm.Error.Printf("unable to encode bookmarks: %v\n", err)
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		pterm.Error.Printf("unable to create destination directory: %v\n", err)
		return err
	}

	if err := os.WriteFile(dest, data, 0o644); err != nil {
		pterm.Error.Printf("unable to write export file: %v\n", err)
		return err
	}

	pterm.Success.Printf("exported %d bookmarks to %s\n", len(bookmarks), dest)
	return nil
}

func init() {
	rootCmd.AddCommand(ExportBookmarksCmd)
}
