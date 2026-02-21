package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var browserImport bool

var ImportBookmarksCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Import bookmarks from a file, replacing the current set",
	RunE:  importBookmarks,
}

var htmlLinkRe = regexp.MustCompile(`(?i)<A HREF="([^"]+)"[^>]*>([^<]+)</A>`)

func importBookmarks(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		pterm.Error.Println("bm import <path>")
		return errors.New("wrong usage")
	}

	src := args[0]

	data, err := os.ReadFile(src)
	if err != nil {
		pterm.Error.Printf("unable to read %s: %v\n", src, err)
		return err
	}

	if browserImport {
		return importBrowserBookmarks(data, src)
	}

	var bookmarks []internal.Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		pterm.Error.Printf("invalid bookmark file: %v\n", err)
		return err
	}

	if err := internal.WriteBookmarks(&bookmarks); err != nil {
		pterm.Error.Printf("unable to save imported bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("imported %d bookmarks from %s\n", len(bookmarks), src)
	return nil
}

func importBrowserBookmarks(data []byte, src string) error {
	existing, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	existingTargets := make(map[string]bool)
	for _, b := range existing {
		existingTargets[b.Target] = true
	}

	matches := htmlLinkRe.FindAllSubmatch(data, -1)
	if len(matches) == 0 {
		pterm.Warning.Println("no bookmarks found in HTML file")
		return nil
	}

	added := 0
	skipped := 0
	for _, m := range matches {
		href := string(m[1])
		title := string(m[2])

		if existingTargets[href] {
			skipped++
			continue
		}

		name := uniqueName(existing, sanitizeName(title))

		existing = append(existing, internal.Bookmark{
			Name:      name,
			Target:    href,
			Type:      "url",
			Tags:      []string{},
			CreatedAt: time.Now().Format(time.RFC3339),
		})
		existingTargets[href] = true
		added++
	}

	if err := internal.WriteBookmarks(&existing); err != nil {
		pterm.Error.Printf("unable to save imported bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("imported %d bookmarks from %s (%d skipped as duplicates)\n", added, src, skipped)
	return nil
}

func sanitizeName(title string) string {
	re := regexp.MustCompile(`\s+`)
	name := re.ReplaceAllString(title, "_")
	runes := []rune(name)
	if len(runes) > 40 {
		name = string(runes[:40])
	}
	return name
}

func uniqueName(bookmarks []internal.Bookmark, base string) string {
	existing := make(map[string]bool)
	for _, b := range bookmarks {
		existing[b.Name] = true
	}
	if !existing[base] {
		return base
	}
	for i := 2; ; i++ {
		candidate := pterm.Sprintf("%s_%d", base, i)
		if !existing[candidate] {
			return candidate
		}
	}
}

func init() {
	ImportBookmarksCmd.Flags().BoolVar(&browserImport, "browser", false, "import from Netscape HTML bookmark format (browser export)")
	rootCmd.AddCommand(ImportBookmarksCmd)
}
