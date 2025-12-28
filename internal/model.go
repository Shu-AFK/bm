package internal

type Bookmark struct {
	Name      string   `json:"name"`
	Target    string   `json:"target"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	Args      []string `json:"args,omitempty"`
	CreatedAt string   `json:"created_at"`
}
