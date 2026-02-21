package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate <path>",
	Short: "Convert an old bookmarks.json to the current format and save it",
	Long: `Reads an old bookmarks.json file, fills in default values for any fields
introduced in newer versions (note, pinned, last_accessed, access_count),
and writes the result to the current bm storage, replacing existing bookmarks.`,
	RunE: migrate,
}

func migrate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("usage: bm migrate <path>")
		return errors.New("wrong usage")
	}

	src := args[0]
	data, err := os.ReadFile(src)
	if err != nil {
		pterm.Error.Printf("unable to read %s: %v\n", src, err)
		return err
	}

	var raw []internal.Bookmark
	if err := json.Unmarshal(data, &raw); err != nil {
		pterm.Error.Printf("invalid bookmark file: %v\n", err)
		return err
	}

	now := time.Now().Format(time.RFC3339)
	for i := range raw {
		if raw[i].Tags == nil {
			raw[i].Tags = []string{}
		}
		if raw[i].Args == nil {
			raw[i].Args = []string{}
		}
		if raw[i].CreatedAt == "" {
			raw[i].CreatedAt = now
		}
		if raw[i].Type == "" {
			raw[i].Type = "url"
		}
	}

	if err := internal.WriteBookmarks(&raw); err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("migrated %d bookmarks from %s\n", len(raw), src)
	return nil
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
