package internal

type Bookmark struct {
	Name         string   `json:"name"`
	Target       string   `json:"target"`
	Type         string   `json:"type"`
	Tags         []string `json:"tags"`
	Args         []string `json:"args,omitempty"`
	CreatedAt    string   `json:"created_at"`
	Note         string   `json:"note,omitempty"`
	Pinned       bool     `json:"pinned,omitempty"`
	LastAccessed string   `json:"last_accessed,omitempty"`
	AccessCount  int      `json:"access_count,omitempty"`
}

type Group struct {
	Name      string   `json:"name"`
	Bookmarks []string `json:"bookmarks"`
	CreatedAt string   `json:"created_at"`
}
