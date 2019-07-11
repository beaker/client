package api

import "time"

// Workspace is the consumabable information about a workspace
type Workspace struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Size        int    `json:"size"` // Total count of items in the workspace.

	Owner  Identity `json:"owner"`
	Author Identity `json:"author"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Archived bool      `json:"archived"`
}

// WorkspacePage is a page of results from a batch workspace API.
type WorkspacePage struct {
	// Results of a batch query.
	Data []Workspace `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}
