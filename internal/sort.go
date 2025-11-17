package internal

import (
	"cmp"
	"slices"
	"time"
)

type SortMode int

const (
	SortByName SortMode = iota
	SortByCreated
	SortByTarget
	SortByType
)

func SortBookmarks(bookmarks []Bookmark, mode SortMode) {
	switch mode {
	case SortByName:
		slices.SortFunc(bookmarks, func(a, b Bookmark) int {
			return cmp.Compare(a.Name, b.Name)
		})

	case SortByCreated:
		slices.SortFunc(bookmarks, func(a, b Bookmark) int {
			ta, errA := time.Parse(time.RFC3339, a.CreatedAt)
			tb, errB := time.Parse(time.RFC3339, b.CreatedAt)

			if errA != nil || errB != nil {
				return cmp.Compare(a.CreatedAt, b.CreatedAt)
			}

			switch {
			case ta.Before(tb):
				return -1
			case ta.After(tb):
				return 1
			default:
				return 0
			}
		})

	case SortByTarget:
		slices.SortFunc(bookmarks, func(a, b Bookmark) int {
			return cmp.Compare(a.Target, b.Target)
		})

	case SortByType:
		slices.SortFunc(bookmarks, func(a, b Bookmark) int {
			return cmp.Compare(a.Type, b.Type)
		})
	}
}
