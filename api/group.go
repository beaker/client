package api

import (
	"path"
	"time"
)

// CreateGroupResponse is a service response returned when a new group is
// created. For now it's just the group ID, but may be expanded in the future.
type CreateGroupResponse struct {
	ID string `json:"id"`
}

// GroupSpec is a specification for creating a new Group.
type GroupSpec struct {
	// (required) Workspace where this group should be placed.
	Workspace string `json:"workspace,omitempty"`

	// (required) Unique name to assign the group.
	Name string `json:"name"`

	// (optional) Text description for the dataset.
	Description string `json:"description,omitempty"`

	// (optional) Initial set of experiments to add to the group.
	Experiments []string `json:"experiments,omitempty"`

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty"`
}

// Group is a collection of experiments.
type Group struct {
	// Identity
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`

	// Ownership
	Owner     Identity           `json:"owner"` // TODO: Deprecated. Refer to containing workspace instead.
	Author    Identity           `json:"author"`
	Workspace WorkspaceReference `json:"workspaceRef"` // TODO: Rename to "workspace" when clients are updated.

	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Archived    bool      `json:"archived"`
}

// DisplayID returns the most human-friendly name available for a group
// while guaranteeing that it's unique and non-empty.
func (e *Group) DisplayID() string {
	if e.Name != "" {
		return path.Join(e.Owner.Name, e.Name)
	}
	return e.ID
}

// GroupExperimentTask identifies an (experiment, task) pair within a group.
type GroupExperimentTask struct {
	Experiment GroupExperiment `json:"experiment"`
	Task       GroupTask       `json:"task"`
}

// GroupExperiment is a minimal experiment summary for aggregated views.
type GroupExperiment struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// GroupTask is a minimal task summary for aggregated views.
type GroupTask struct {
	ID        string                 `json:"id"`
	LastState *ExecutionState        `json:"lastState,omitempty"`
	Canceled  *time.Time             `json:"canceled,omitempty"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	Env       map[string]string      `json:"env,omitempty"`
	Name      string                 `json:"name,omitempty"`
}

// GroupPatchSpec describes a patch to apply to a group's editable fields.
type GroupPatchSpec struct {
	// (optional) Unqualified name to assign to the group. It is considered
	// a collision error if another group has the same creator and name.
	Name *string `json:"name,omitempty"`

	// (optional) Description to assign to the group or empty string to
	// delete an existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Experiment IDs to add to the group.
	// It is an error to add and remove the same experiment in one patch.
	AddExperiments []string `json:"addExperiments,omitempty"`

	// (optional) Experiment IDs to remove from the group.
	// It is an error to add and remove the same experiment in one patch.
	RemoveExperiments []string `json:"removeExperiments,omitempty"`

	// (optional) New selected environment variables and metrics.
	Parameters *[]GroupParameter `json:"parameters,omitempty"`

	// (optional) Whether the group should be archived. Ignored if nil.
	Archive *bool `json:"archive,omitempty"`
}

// GroupParameterType enumerates sources for a paramater in group analyses.
type GroupParameterType string

const (
	// MetricParameter is a parameter parsed from a user-generated metric file in a task result.
	MetricParameter GroupParameterType = "metric"

	// EnvVarParameter is a parameter supplied as an environment variable in a task spec.
	EnvVarParameter GroupParameterType = "env"
)

// GroupParameter is a measurable value for use in group analyses.
type GroupParameter struct {
	Type GroupParameterType `json:"type"`
	Name string             `json:"name"`
}

// GroupParameterCount summarizes how often a parameter is observed among a group's tasks.
type GroupParameterCount struct {
	Type  GroupParameterType `json:"type"`
	Name  string             `json:"name"`
	Count int64              `json:"count"`
}

// GroupPage is a page of results from a batch group API.
type GroupPage struct {
	// Results of a batch query.
	Data []Group `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}
