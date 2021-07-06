package api

import (
	"time"

	"github.com/allenai/bytefmt"
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
	ID         string `json:"id"`
	Task       string `json:"task"`
	Experiment string `json:"experiment"`
	Workspace  string `json:"workspace"`

	// Node is set when a task has been assigned to a node.
	Node string `json:"node,omitempty"`

	// Spec describes the execution's task, but with all soft references fully resolved.
	Spec   TaskSpecV2   `json:"spec"`
	Result ResultTarget `json:"result"`

	// State describes execution status and progression.
	State ExecutionState `json:"state"`

	// Limits describes resources assigned to this execution
	Limits ResourceLimits `json:"limits"`

	// (deprecated) See corresponding value in Spec.
	Priority Priority `json:"priority"`
}

// ResultTarget describes a target to which results will be written.
type ResultTarget struct {
	// Name or ID of a Beaker dataset.
	Beaker string `json:"beaker,omitempty"`
}

// ResourceLimits describes limits assigned to a process.
type ResourceLimits struct {
	// (optional) CPUCount sets a minimum number of logical CPU cores and may be fractional.
	//
	// Examples: 4, 0.5
	CPUCount float64 `json:"cpuCount,omitempty"`

	// (optional) GPUs assigned to the session. Either GPU index or ID.
	GPUs []string `json:"gpus,omitempty"`

	// (optional) Memory sets a limit for CPU memory, which may be a raw number
	// of bytes or a formatted string with a number followed by a unit suffix.
	//
	// Examples: "2.5 GiB", 2684354560
	Memory *bytefmt.Size `json:"memory,omitempty"`
}

// ExecutionState details an execution's status.
type ExecutionState struct {
	Created   time.Time  `json:"created"`
	Scheduled *time.Time `json:"scheduled,omitempty"`
	Started   *time.Time `json:"started,omitempty"`
	Exited    *time.Time `json:"exited,omitempty"`
	Failed    *time.Time `json:"failed,omitempty"`
	Finalized *time.Time `json:"finalized,omitempty"`

	// ExitCode is an integer process exit code, if the process exited normally.
	ExitCode *int `json:"exitCode,omitempty"`

	// Message describes additional state-related context.
	Message string `json:"message,omitempty"`

	// Canceled indicates whether and when an execution was canceled.
	Canceled *time.Time `json:"canceled,omitempty"`
}

// ExecStatusUpdate changes a process' status. Unset fields are left unchanged.
// Timestamp fields can only be written once and will be ignored if already set.
type ExecStatusUpdate struct {
	// (optional) Scheduled is set when a process has been assigned to a node.
	Scheduled bool `json:"scheduled,omitempty"`

	// (optional) Started is set when a process has started running.
	Started bool `json:"started,omitempty"`

	// (optional) ExitCode is set when a process exits normally.
	ExitCode *int `json:"exitCode,omitempty"`

	// (optional) Failed is set when a process has ended abnormally.
	Failed bool `json:"failed,omitempty"`

	// (optional) Finalized is set when a process has ended and all results have been captured.
	Finalized bool `json:"finalized,omitempty"`

	// (optional) Canceled is set to terminate a process remotely.
	Canceled bool `json:"canceled,omitempty"`

	// (optional) Human-readable message to provide context for the status.
	Message *string `json:"message,omitempty"`

	// (optional) Limits record the maximum resources available during execution.
	Limits *ResourceLimits `json:"limits,omitempty"`
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
