package internal

import "slices"

func SearchByTag(bookmarks []Bookmark, tag string) []Bookmark {
	var ret []Bookmark
	for _, bm := range bookmarks {
		if slices.Contains(bm.Tags, tag) {
			ret = append(ret, bm)
		}
	}

	return ret
}
