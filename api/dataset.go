package api

import (
	"time"
)

// DatasetStorage is a reference to a FileHeap dataset.
type DatasetStorage struct {
	Address      string    `json:"address"`
	ID           string    `json:"id"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"tokenExpires"`
}

// Dataset is a file or collection of files. It may be the result of a task or
// uploaded directly by a user.
type Dataset struct {
	// Identity
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"fullName,omitempty"`

	// Ownership
	Owner     Identity           `json:"owner"` // TODO: Deprecated. Refer to containing workspace instead.
	Author    Identity           `json:"author"`
	Workspace WorkspaceReference `json:"workspaceRef"` // TODO: Rename to "workspace" when clients are updated.

	// Status
	Created   time.Time `json:"created"`
	Committed time.Time `json:"committed,omitempty"`
	Archived  bool      `json:"archived"`

	// A plain-text description of this dataset.
	Description string `json:"description,omitempty"`

	// Task for which this dataset is a result, i.e. provenance, if any.
	SourceTask *string `json:"sourceTask,omitempty"`

	Storage *DatasetStorage `json:"storage,omitempty"`
}

// DatasetSpec is a specification for creating a new Dataset.
type DatasetSpec struct {
	// (required) Workspace where this dataset should be placed.
	Workspace string `json:"workspace,omitempty"`

	// (optional) Text description for the dataset.
	Description string `json:"description,omitempty"`

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty"`

	// (optional) If set, the dataset will be stored in FileHeap.
	// This flag will eventually become the default and be removed.
	FileHeap bool `json:"fileHeap,omitempty"`
}

// DatasetPatchSpec describes a patch to apply to a dataset's editable fields.
// Only one field may be set in a single request.
type DatasetPatchSpec struct {
	// (optional) Unqualified name to assign to the dataset. It is considered
	// a collision error if another dataset has the same creator and name.
	Name *string `json:"name,omitempty"`

	// (optional) Description to assign to the dataset or empty string to
	// delete an existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Whether the dataset should be locked for writes. Ignored if false.
	Commit bool `json:"commit,omitempty"`

	// (optional) Whether the dataset should be archived. Ignored if nil.
	Archive *bool `json:"archive,omitempty"`
}

// DatasetPage is a page of results from a batch dataset API.
type DatasetPage struct {
	// Results of a batch query.
	Data []Dataset `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}
