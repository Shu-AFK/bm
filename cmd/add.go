package cmd

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Shu-AFK/bm/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var tags []string
var unreachableErr = errors.New("unreachable type")

var AddBookmarkCmd = &cobra.Command{
	Use:   "add <name> <path/url>",
	Short: "Add a new bookmark",
	RunE:  add,
}

func normalizeURL(s string) (string, bool) {
	u, err := url.Parse(s)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		return u.String(), true
	}
	if strings.Contains(s, ".") && !strings.ContainsAny(s, `/\ `) {
		return "https://" + s, true
	}
	return "", false
}

func resolvePath(p string) (string, os.FileInfo, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", nil, err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", nil, errors.New("path does not exist")
	}
	fi, err := os.Stat(real)
	if err != nil {
		return "", nil, errors.New("path does not exist")
	}
	return real, fi, nil
}

func getBookmarkInfo(target string) (string, string, error) {
	if u, ok := normalizeURL(target); ok {
		return u, "url", nil
	}
	real, fi, err := resolvePath(target)
	if err != nil {
		return "", "", err
	}
	if fi.IsDir() {
		return real, "folder", nil
	}
	if fi.Mode().IsRegular() {
		return real, "file", nil
	}
	return "", "", unreachableErr
}

func add(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		pterm.Error.Println("usage: bm add <name> <path/url>")
		return errors.New("wrong usage")
	}
	name := args[0]
	target := args[1]

	bookmarks, err := internal.ReadBookmarks()
	if err != nil {
		pterm.Error.Printf("unable to read bookmarks: %v\n", err)
		return err
	}

	for _, b := range bookmarks {
		if b.Name == name {
			pterm.Error.Println("a bookmark with this name already exists")
			return errors.New("name already exists")
		}
	}

	resolved, bmType, err := getBookmarkInfo(target)
	if err != nil {
		pterm.Error.Printf("invalid path/url: %v\n", err)
		return err
	}

	bookmarks = append(bookmarks, internal.Bookmark{
		Name:      name,
		Target:    resolved,
		Type:      bmType,
		Tags:      tags,
		CreatedAt: time.Now().String(),
	})

	err = internal.WriteBookmarks(&bookmarks)
	if err != nil {
		pterm.Error.Printf("unable to write bookmarks: %v\n", err)
		return err
	}

	pterm.Success.Printf("added %s (%s)\n", name, bmType)
	return nil
}

func init() {
	AddBookmarkCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "All tags associated with that bookmark, comma seperated")
	rootCmd.AddCommand(AddBookmarkCmd)
}
