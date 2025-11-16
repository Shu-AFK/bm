package internal

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type MatchKind int

const (
	MatchNone MatchKind = iota
	MatchExact
	MatchPrefix
	MatchFuzzy
)

const fuzzyThreshold = 5

type MatchResult struct {
	Kind       MatchKind
	BM         *Bookmark
	Candidates []*Bookmark
}

func SearchByTag(bookmarks []Bookmark, tag string) []Bookmark {
	tagLower := strings.ToLower(tag)
	var ret []Bookmark

	for _, bm := range bookmarks {
		for _, t := range bm.Tags {
			if strings.ToLower(t) == tagLower {
				ret = append(ret, bm)
				break
			}
		}
	}

	return ret
}

func GetClosestMatch(bookmarks []Bookmark, query string) MatchResult {
	q := strings.ToLower(query)

	// 1) Exact match
	for i := range bookmarks {
		bm := &bookmarks[i]
		if strings.ToLower(bm.Name) == q {
			return MatchResult{Kind: MatchExact, BM: bm}
		}
	}

	// 2) Prefix match
	var prefixes []*Bookmark
	for i := range bookmarks {
		if strings.HasPrefix(strings.ToLower(bookmarks[i].Name), q) {
			prefixes = append(prefixes, &bookmarks[i])
		}
	}

	if len(prefixes) == 1 {
		return MatchResult{Kind: MatchPrefix, BM: prefixes[0]}
	}
	if len(prefixes) > 1 {
		return MatchResult{Kind: MatchPrefix, Candidates: prefixes}
	}

	// 3) Fuzzy Match
	names := make([]string, len(bookmarks))
	for i, b := range bookmarks {
		names[i] = strings.ToLower(b.Name)
	}

	matches := fuzzy.RankFind(q, names)
	filtered := fuzzyFilter(matches)
	if len(filtered) == 0 {
		return MatchResult{Kind: MatchNone}
	}
	if len(filtered) == 1 {
		return MatchResult{Kind: MatchFuzzy, BM: &bookmarks[filtered[0].OriginalIndex]}
	}
	sort.Sort(filtered)

	fuzzyBm := make([]*Bookmark, len(filtered))
	for i := range filtered {
		fuzzyBm[i] = &bookmarks[filtered[i].OriginalIndex]
	}
	return MatchResult{Kind: MatchFuzzy, Candidates: fuzzyBm}
}

func fuzzyFilter(matches fuzzy.Ranks) fuzzy.Ranks {
	var filtered fuzzy.Ranks
	for _, m := range matches {
		if m.Distance <= fuzzyThreshold {
			filtered = append(filtered, m)
		}
	}

	return filtered
}
