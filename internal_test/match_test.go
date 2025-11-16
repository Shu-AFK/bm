package internal_test

import (
	"testing"

	"github.com/Shu-AFK/bm/internal"
)

func TestSearchByTag_CaseInsensitive(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "Go", Tags: []string{"Lang", "CLI"}},
		{Name: "Rust", Tags: []string{"lang", "System"}},
		{Name: "Python", Tags: []string{"scripting"}},
	}

	res := internal.SearchByTag(bookmarks, "LANG")

	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}

	names := []string{res[0].Name, res[1].Name}
	if names[0] != "Go" && names[1] != "Rust" {
		t.Fatalf("expected Go + Rust, got %+v", names)
	}
}

func TestGetClosestMatch_Exact(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "Google"},
		{Name: "GitHub"},
		{Name: "Gmail"},
	}

	res := internal.GetClosestMatch(bookmarks, "GitHub")

	if res.Kind != internal.MatchExact {
		t.Fatalf("expected MatchExact, got %v", res.Kind)
	}
	if res.BM == nil || res.BM.Name != "GitHub" {
		t.Fatalf("expected GitHub, got %+v", res.BM)
	}
}

func TestGetClosestMatch_PrefixSingle(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "go"},
		{Name: "github"},
		{Name: "gitlab"},
	}

	res := internal.GetClosestMatch(bookmarks, "gith")

	if res.Kind != internal.MatchPrefix {
		t.Fatalf("expected MatchPrefix, got %v", res.Kind)
	}
	if res.BM == nil || res.BM.Name != "github" {
		t.Fatalf("expected github, got %+v", res.BM)
	}
}

func TestGetClosestMatch_PrefixMultiple(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "git"},
		{Name: "github"},
		{Name: "gitlab"},
	}

	res := internal.GetClosestMatch(bookmarks, "gi")

	if res.Kind != internal.MatchPrefix {
		t.Fatalf("expected MatchPrefix, got %v", res.Kind)
	}
	if len(res.Candidates) != 3 {
		t.Fatalf("expected 3 candidates, got %d", len(res.Candidates))
	}
}

func TestGetClosestMatch_FuzzySingle(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "google"},
		{Name: "github"},
		{Name: "gitlab"},
	}

	res := internal.GetClosestMatch(bookmarks, "gogle") // typo

	if res.Kind != internal.MatchFuzzy {
		t.Fatalf("expected MatchFuzzy, got %v", res.Kind)
	}
	if res.BM == nil || res.BM.Name != "google" {
		t.Fatalf("expected fuzzy match google, got %+v", res.BM)
	}
}

func TestGetClosestMatch_FuzzyMultiple(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "google"},
		{Name: "gogole"},
		{Name: "goggle"},
	}

	res := internal.GetClosestMatch(bookmarks, "gogle")

	if res.Kind != internal.MatchFuzzy {
		t.Fatalf("expected MatchFuzzy, got %v", res.Kind)
	}
	if len(res.Candidates) < 2 {
		t.Fatalf("expected >=2 fuzzy candidates, got %d", len(res.Candidates))
	}
}

func TestGetClosestMatch_None(t *testing.T) {
	bookmarks := []internal.Bookmark{
		{Name: "google"},
		{Name: "github"},
	}

	res := internal.GetClosestMatch(bookmarks, "xxxxxxxx")

	if res.Kind != internal.MatchNone {
		t.Fatalf("expected MatchNone, got %v", res.Kind)
	}
	if res.BM != nil {
		t.Fatalf("BM should be nil for MatchNone")
	}
	if len(res.Candidates) != 0 {
		t.Fatalf("expected 0 candidates, got %d", len(res.Candidates))
	}
}
