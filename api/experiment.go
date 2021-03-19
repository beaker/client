package api

import (
	"time"
)

// Experiment describes an experiment and its tasks.
type Experiment struct {
	// Identity
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"fullName,omitempty"`

	// Ownership
	Owner     Identity           `json:"owner"` // TODO: Deprecated. Refer to containing workspace instead.
	Author    Identity           `json:"author"`
	Workspace WorkspaceReference `json:"workspaceRef"` // TODO: Rename to "workspace" when clients are updated.

	Description string       `json:"description,omitempty"`
	Executions  []*Execution `json:"executions,omitempty"`
	Created     time.Time    `json:"created"`
	Canceled    bool         `json:"canceled,omitempty"`
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
}

// ExperimentPage is a page of results from a batch experiment API.
type ExperimentPage struct {
	// Results of a batch query.
	Data []Experiment `json:"data"`

	// Opaque token to the element after Data, provided only if more data is available.
	NextCursor string `json:"nextCursor,omitempty"`
}

// Task is a full description of a task specification and its status.
type Task struct {
	// Identity
	ID           string `json:"id"`
	ExperimentID string `json:"experimentId"`
	Name         string `json:"name,omitempty"`

	// Ownership
	Owner  Identity `json:"owner"`
	Author Identity `json:"author"`

	// State of this task and its execution(s).
	Created     time.Time   `json:"created"`
	Schedulable bool        `json:"schedulable"`
	Executions  []Execution `json:"executions,omitempty"`
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

	// Canceled indicates whether and when an execution was canceled.
	Canceled *time.Time `json:"canceled,omitempty"`
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

// ExecutionPatchSpec describes a patch to apply to a execution's editable fields.
type ExecutionPatchSpec struct {
	// (optional) Priority to assign to the execution.
	Priority Priority `json:"priority,omitempty"`
}

// ExecutionResults is the structured content of an execution's output 'metrics.json' file.
type ExecutionResults struct {
	Metrics map[string]interface{} `json:"metrics"`
}
