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
	Short: "Migrate bookmarks from an old JSON file to the current storage",
	Long: `Reads an old bookmarks.json file, fills in default values for any fields
introduced in newer versions, and writes the result to the current bm storage.

By default the current bookmarks are replaced. Use --merge to add only
bookmarks whose names do not already exist in the current storage.`,
	RunE: migrate,
}

var migrateMerge bool

// oldBookmark accepts every field that has ever existed so unknown fields
// are preserved rather than silently dropped.
type oldBookmark struct {
	Name         string   `json:"name"`
	Target       string   `json:"target"`
	Type         string   `json:"type"`
	Tags         []string `json:"tags"`
	Args         []string `json:"args"`
	CreatedAt    string   `json:"created_at"`
	Note         string   `json:"note"`
	Pinned       bool     `json:"pinned"`
	LastAccessed string   `json:"last_accessed"`
	AccessCount  int      `json:"access_count"`
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

	var raw []oldBookmark
	if err := json.Unmarshal(data, &raw); err != nil {
		pterm.Error.Printf("invalid bookmark file: %v\n", err)
		return err
	}

	now := time.Now().Format(time.RFC3339)
	imported := make([]internal.Bookmark, 0, len(raw))
	for _, r := range raw {
		bm := internal.Bookmark{
			Name:         r.Name,
			Target:       r.Target,
			Type:         r.Type,
			Tags:         r.Tags,
			Args:         r.Args,
			CreatedAt:    r.CreatedAt,
			Note:         r.Note,
			Pinned:       r.Pinned,
			LastAccessed: r.LastAccessed,
			AccessCount:  r.AccessCount,
		}

		if bm.Tags == nil {
			bm.Tags = []string{}
		}
		if bm.Args == nil {
			bm.Args = []string{}
		}
		if bm.CreatedAt == "" {
			bm.CreatedAt = now
		}
		if bm.Type == "" {
			bm.Type = "url"
		}

		imported = append(imported, bm)
	}

	if !migrateMerge {
		if err := internal.WriteBookmarks(&imported); err != nil {
			pterm.Error.Printf("unable to write bookmarks: %v\n", err)
			return err
		}
		pterm.Success.Printf("migrated %d bookmarks from %s\n", len(imported), src)
		return nil
	}

	existing, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read existing bookmarks: %v\n", err)
		return err
	}

	existingNames := make(map[string]bool)
	for _, b := range existing {
		existingNames[b.Name] = true
	}

	added := 0
	skipped := 0
	for _, bm := range imported {
		if existingNames[bm.Name] {
			skipped++
			continue
		}
		existing = append(existing, bm)
		existingNames[bm.Name] = true
		added++
	}

	if err := internal.WriteBookmarks(&existing); err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("merged %d bookmarks from %s (%d skipped, name already exists)\n", added, src, skipped)
	return nil
}

func init() {
	migrateCmd.Flags().BoolVar(&migrateMerge, "merge", false, "merge into existing bookmarks instead of replacing them")
	rootCmd.AddCommand(migrateCmd)
}
