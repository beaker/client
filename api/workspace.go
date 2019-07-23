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

// WorkspaceSpec is a specification for creating a new Workspace
type WorkspaceSpec struct {
	// (required) Name of the workspace to be created
	Name string `json:"name"`

	// (optional) Description of the workspace to be created
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public"`

	// (required) Organization on behalf of whom this resource is created. The
	// user issuing the request must be a member of the organization.
	Owner string `json:"owner,omitempty"`
}

// WorkspacePage is a page of results from a batch workspace API.
type WorkspacePage struct {
	// Results of a batch query.
	Data []Workspace `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
	// (required) Organization on whose behalf this resource is created. The
	// user issuing the request must be a member of the organization.
	Organization string `json:"org"`
}

// CreateWorkspaceResponse is a service response returned when a new workspace is
// created.
type CreateWorkspaceResponse struct {
	ID string `json:"id"`
}
