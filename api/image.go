package api

import (
	"time"
)

// Image describes the Docker image ran by a Task while executing an Experiment.
type Image struct {
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

	// Original image tag, if supplied on creation. See ImageSpec.
	OriginalTag string `json:"originalTag,omitempty"`

	// A plain-text description of this image.
	Description string `json:"description,omitempty"`
}

// ImageSpec is a specification for creating a new Image.
type ImageSpec struct {
	// (required) Workspace where this image should be placed.
	Workspace string `json:"workspace,omitempty"`

	// (required) Unique identifier for the image's image. In Docker images,
	// this is a SHA256 hash.
	ImageID string `json:"ImageID"` // TODO: convert to loweCase name and update reflection tag

	// (optional) Text description for the image.
	Description string `json:"Description,omitempty"` // TODO: convert to loweCase name and update reflection tag

	// (optional) Original image tag from which the image was created.
	ImageTag string `json:"ImageTag,omitempty"` // TODO: convert to loweCase name and update reflection tag

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty"`
}

// ImagePatchSpec describes a patch to apply to an image's editable fields.
// Only one field may be set in a single request.
type ImagePatchSpec struct {
	// (optional) Unqualified name to assign to the image. It is considered
	// a collision error if another image has the same creator and name.
	Name *string `json:"name,omitempty"`

	// (optional) Description to assign to the image or empty string to
	// delete an existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Whether the image should be committed. Ignored if false.
	// When committed, an image is placed in Beaker's internal registry and
	// further attempts to push its image will be ignored.
	Commit bool `json:"commit,omitempty"`
}

// ImageRepository contains a repository/tag and credentials required to
// upload an image's Docker image via "docker push".
type ImageRepository struct {
	// Full tag, including registry, expected by Beaker. Clients must push this
	// tag exactly for Beaker to recognize the image.
	ImageTag string `json:"imageTag"`

	// Credentials for the image's registry.
	Auth RegistryAuth `json:"auth"`
}

// RegistryAuth supplies authorization for a Docker registry.
type RegistryAuth struct {
	ServerAddress string `json:"server_address"`
	User          string `json:"user"`
	Password      string `json:"password"`
}

// ImagePage is a page of results from a batch image API.
type ImagePage struct {
	// Results of a batch query.
	Data []Image `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}
