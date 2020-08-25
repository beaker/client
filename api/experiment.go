package api

import (
	"path"
	"time"
)

// Experiment describes an experiment and its tasks.
type Experiment struct {
	// Identity
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`

	// Ownership
	Owner     Identity           `json:"owner"` // TODO: Deprecated. Refer to containing workspace instead.
	Author    Identity           `json:"author"`
	Workspace WorkspaceReference `json:"workspaceRef"` // TODO: Rename to "workspace" when clients are updated.

	Description string           `json:"description,omitempty"`
	Nodes       []ExperimentNode `json:"nodes"`
	Created     time.Time        `json:"created"`
	Archived    bool             `json:"archived"`
}

// DisplayID returns the most human-friendly name available for an experiment
// while guaranteeing that it's unique and non-empty.
func (e *Experiment) DisplayID() string {
	if e.Name != "" {
		return path.Join(e.Owner.Name, e.Name)
	}
	return e.ID
}

// ExperimentSpec describes a set of tasks with optional dependencies.
// This set represents a (potentially disconnected) directed acyclic graph.
type ExperimentSpec struct {
	// (optional) Version must be 'v1' or left unset.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// (required) Workspace where this experiment and its results should be placed.
	Workspace string `json:"workspace,omitempty" yaml:"workspace,omitempty"`

	// (optional) Text description of the experiment.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// (required) Tasks to create. Tasks may be defined in any order, though all
	// dependencies must be internally resolvable within the experiment.
	Tasks []ExperimentTaskSpec `json:"tasks" yaml:"tasks"`

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty" yaml:"authorToken,omitempty"`
}

// ExperimentNode describes a task along with its links within an experiment.
type ExperimentNode struct {
	Name      string          `json:"name,omitempty"`
	TaskID    string          `json:"taskId"`
	ResultID  string          `json:"resultId"`
	LastState *ExecutionState `json:"lastState,omitempty"`
	Canceled  *time.Time      `json:"canceled,omitempty"`
	Status    TaskStatus      `json:"status"`

	// Identifiers of tasks dependent on this node within the containing experiment.
	ChildTasks []string `json:"childTaskIds"`

	// Identifiers of task on which this node depends within the containing experiment.
	ParentTasks []string `json:"parentTaskIds"`
}

// DisplayID returns the most human-friendly name available for an experiment
// node while guaranteeing that it's unique within the context of its experiment.
func (n *ExperimentNode) DisplayID() string {
	if n.Name != "" {
		return n.Name
	}
	return n.TaskID
}

// ExperimentTaskSpec describes a task spec with optional dependencies on other
// tasks within an experiment. Tasks refer to each other by the Name field.
type ExperimentTaskSpec struct {
	// (optional) Name of the task node, which need only be defined if
	// dependencies reference it.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// (required) Specification describing the task to run.
	Spec TaskSpec `json:"spec" yaml:"spec,omitempty"`

	// (optional) Tasks on which this task depends. Mounts will be applied, in
	// the order defined here, after existing mounts in the task spec.
	DependsOn []TaskDependency `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`

	// (optional) Name of a cluster on which the task should run.
	// Cluster affinity supercedes task requirements.
	Cluster string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

// TaskDependency describes a single "edge" in a task dependency graph.
type TaskDependency struct {
	// (required) Name of the task on which the referencing task depends.
	ParentName string `json:"parentName" yaml:"parentName"`

	// (optional) Path in the child task to which parent results will be mounted.
	// If absent, this is treated as an order-only dependency.
	ContainerPath string `json:"containerPath,omitempty" yaml:"containerPath,omitempty"`
}

// ExperimentPatchSpec describes a patch to apply to an experiment's editable
// fields. Only one field may be set in a single request.
type ExperimentPatchSpec struct {
	// (optional) Unqualified name to assign to the experiment. It is considered
	// a collision error if another experiment has the same creator and name.
	Name *string `json:"name,omitempty"`

	// (optional) Description to assign to the experiment or empty string to
	// delete an existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Whether the experiment should be archived. Ignored if nil.
	Archive *bool `json:"archive,omitempty"`
}

// ExperimentPage is a page of results from a batch experiment API.
type ExperimentPage struct {
	// Results of a batch query.
	Data []Experiment `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}

// Executions is an ordered collection of task executions.
type Executions struct {
	Data []Execution `json:"data"`
}

// Execution represents an attempt to run a task. A task may have many executions.
type Execution struct {
	ID   string `json:"id"`
	Task string `json:"task"`

	// Node is set when a task has been assigned to a node.
	Node string `json:"node,omitempty"`

	// Spec describes the execution's task, but with all soft references fully resolved.
	Spec   TaskSpecV2   `json:"spec"`
	Result ResultTarget `json:"result"`

	// State describes execution status and progression.
	State ExecutionState `json:"state"`

	// Limits describes resources assigned to this execution
	Limits TaskResources `json:"limits"`

	// (deprecated) See corresponding value in Spec.
	Priority Priority `json:"priority"`
}

// ResultTarget describes a target to which results will be written.
type ResultTarget struct {
	// Name or ID of a Beaker dataset.
	Beaker string `json:"beaker,omitempty"`
}

// ExecutionState details an execution's status.
type ExecutionState struct {
	Created   time.Time  `json:"created"`
	Scheduled *time.Time `json:"scheduled,omitempty"`
	Started   *time.Time `json:"started,omitempty"`
	Ended     *time.Time `json:"ended,omitempty"`
	Finalized *time.Time `json:"finalized,omitempty"`

	// ExitCode is an integer process exit code, if the process exited normally.
	ExitCode *int `json:"exitCode,omitempty"`

	// Message describes additional state-related context.
	Message string `json:"message,omitempty"`

	// TaskCanceled indicates whether and when an execution's task was canceled.
	TaskCanceled *time.Time `json:"taskCanceled,omitempty"`
}

// ExecStatus describes what phase an execution is in.
type ExecStatus string

const (
	// ExecScheduled indicates a task has been assigned to a node and is preparing to run.
	ExecScheduled ExecStatus = "scheduled"

	// ExecStarted indicates an execution has started running.
	ExecStarted ExecStatus = "started"

	// ExecEnded indicates an execution exited or was interrupted. Its results
	// may not yet be final. Callers may inspect its exit code to determine
	// success or failure. If none is set, the execution is considered failed.
	ExecEnded ExecStatus = "ended"

	// ExecFinalized indicates a task has ended and all results have been captured.
	ExecFinalized ExecStatus = "finalized"
)

// ExecStatusUpdate snapshots a task execution's status.
type ExecStatusUpdate struct {
	// (optional) Status is the task's current stage of execution.
	Status ExecStatus `json:"status,omitempty"`

	// (optional) Human-readable message to provide context for the status.
	Message *string `json:"message,omitempty"`

	// (optional) Exit code of the task's process.
	ExitCode *int `json:"exitCode,omitempty"`

	// (optional) Limits record the maximum resources available during execution.
	Limits *TaskResources `json:"limits,omitempty"`
}
