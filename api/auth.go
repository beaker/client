package api

// Permission represent's a user's access to an object.
type Permission string

const (
	// NoPermission indicates a user isn't permitted any access to an object.
	NoPermission Permission = "none"

	// Read allows a user to read an object.
	Read Permission = "read"

	// Write allows a user to read, modify and delete an object.
	Write Permission = "write"

	// FullControl indicates a user can write an object and read or modify its permissions.
	FullControl Permission = "all"
)

// WorkspacePermissionSummary aggregates permissions for a particular workspace. Some fields
// may be omitted when viewed by a user with limited permissions.
// TODO: After migrating permissions to be workspace-based, rename to PermissionSummary
type WorkspacePermissionSummary struct {
	// Authorization for the user issuing a request.
	RequesterAuth Permission `json:"requesterAuth"`

	// Default permissions granted on the object.
	Public bool `json:"public"`

	// Mapping of additional permissions granted to each user, indexed by user ID.
	Authorizations map[string]Permission `json:"authorizations,omitempty"`
}

// WorkspacePermissionPatch describes transactional changes to a single workspace's permissions.
// TODO: After migrating permissions to be workspace-based, rename to PermissionPatch
type WorkspacePermissionPatch struct {
	// (optional) Default permission to apply to all users.
	Public *bool `json:"public,omitempty"`

	// (optional) Mapping of additional permissions granted to each user. Set a
	// user's permission to "none" to reset that user's authorization to default.
	Authorizations map[string]Permission `json:"authorizations,omitempty"`
}
