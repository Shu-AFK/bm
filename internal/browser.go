package internal

import "errors"

func Open(bookmark Bookmark) error {
	var err error

	switch bookmark.Type {
	case "url":
		err = openURL(bookmark.Target)
	case "path":
		err = openPath(bookmark.Target)
	default:
		err = errors.New("unknown bookmark type")
	}

	return err
}

func openURL(url string) error {
	return nil
}

func openPath(path string) error {
	return nil
}
