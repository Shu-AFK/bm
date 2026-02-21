package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func getStoragePath() (string, error) {
	dir := filepath.Join(xdg.DataHome, "bm")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "bookmarks.json"), nil
}

func getGroupsPath() (string, error) {
	dir := filepath.Join(xdg.DataHome, "bm")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "groups.json"), nil
}

func ReadBookmarks() ([]Bookmark, error) {
	path, err := getStoragePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Bookmark{}, nil
		}
		return nil, err
	}

	bookmarks := new([]Bookmark)
	err = json.Unmarshal([]byte(data), bookmarks)
	if err != nil {
		return nil, err
	}

	return *bookmarks, nil
}

func WriteBookmarks(bookmarks *[]Bookmark) error {
	jsonBookmarks, err := json.Marshal(*bookmarks)
	if err != nil {
		return err
	}

	path, err := getStoragePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonBookmarks, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func ReadGroups() ([]Group, error) {
	path, err := getGroupsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Group{}, nil
		}
		return nil, err
	}

	groups := new([]Group)
	err = json.Unmarshal(data, groups)
	if err != nil {
		return nil, err
	}

	return *groups, nil
}

func WriteGroups(groups *[]Group) error {
	jsonGroups, err := json.Marshal(*groups)
	if err != nil {
		return err
	}

	path, err := getGroupsPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonGroups, 0o644)
	if err != nil {
		return err
	}

	return nil
}
