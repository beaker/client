package api

import "time"

// Workspace is the consumabable information about a workspace
type Workspace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	Description string             `json:"description,omitempty"`
	Size        WorkspaceItemCount `json:"size"`

	Owner  Identity `json:"owner"`
	Author Identity `json:"author"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Archived bool      `json:"archived"`
}

// WorkspaceItemCount describes how many items of each type are contained within a workspace.
type WorkspaceItemCount struct {
	Datasets    int `json:"datasets"`
	Experiments int `json:"experiments"`
	Groups      int `json:"groups"`
	Images      int `json:"images"`
}

// WorkspaceSpec is a specification for creating a new Workspace
type WorkspaceSpec struct {
	// (required) Name of the workspace to be created
	Name string `json:"name"`

	// (optional) Description of the workspace to be created
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public"`

	// (optional) Organization on behalf of whom this resource is created. The
	// user issuing the request must be a member of the organization.
	Organization string `json:"org,omitempty"`
}

// WorkspacePatchSpec describes a patch to apply to a workspace's editable fields.
type WorkspacePatchSpec struct {
	// (optional) New name to give the workspace. This will break any existing references by name.
	Name *string `json:"name,omitempty"`

	// (optional) New description to give the workspace
	Description *string `json:"description,omitempty"`

	// (optional) Whether the workspace should be archived. Ignored if nil.
	Archive *bool `json:"archive,omitempty"`
}

type WorkspaceTransferSpec struct {
	IDs []string `json:"ids"`
}

// WorkspacePage is a page of results from a batch workspace API.
type WorkspacePage struct {
	// Results of a batch query.
	Data []Workspace `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`

	// Organization on whose behalf this resource is created. The
	// user issuing the request must be a member of the organization.
	Organization string `json:"org"`
}

// WorkspaceReference is a reference to a workspace in the system, providing both
// name and ID for human-readible and static references, respectively.
type WorkspaceReference struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
}

// Secrets is an ordered collection of secrets.
type Secrets struct {
	Data []Secret `json:"data"`
}

// A Secret describes privileged data stored within a workspace. Secrets can
// only be accessed by workspace contributors with 'write' or higher permission.
type Secret struct {
	Name string `json:"name"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
