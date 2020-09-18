package api

import (
	"time"
)

// Task is a full description of a task specification and its status.
type Task struct {
	// Identity
	ID           string `json:"id"`
	ExperimentID string `json:"experimentId"`
	Name         string `json:"name,omitempty"`

	// Ownership
	Owner  Identity `json:"owner"`
	Author Identity `json:"author"`

	// Status (deprecated fields)
	Status  TaskStatus `json:"status"`
	Created time.Time  `json:"created"`
	Started time.Time  `json:"started"`
	Ended   time.Time  `json:"ended"`

	// State of this task's most recent execution, if any.
	LastState *ExecutionState `json:"lastState,omitempty"`
	Canceled  *time.Time      `json:"canceled,omitempty"`

	// Creation parameters
	Spec        TaskSpec        `json:"spec"`
	ResumedFrom ResumedFromSpec `json:"resumedFrom"`

	// Scheduling
	Cluster string `json:"cluster,omitempty"`
	Node    string `json:"node,omitempty"`

	// Results
	ResultID string `json:"resultId"`
	ExitCode int    `json:"exitCode,omitempty"`
}

type ResumedFromSpec struct {
	TaskID       string `json:"taskId,omitempty"`
	ExperimentID string `json:"experimentId,omitempty"`
}

type TaskLogUploadLink struct {
	TaskID      string `json:"taskId"`
	TaskAttempt string `json:"taskAttempt"`
	LogChunk    string `json:"logChunk"`
	URL         string `json:"url"`
}

type TaskResults struct {
	Metrics map[string]interface{} `json:"metrics"`
}

// TaskSpec contains all information necessary to create a new task.
type TaskSpec struct {
	// (required) Image containing the code to be run.
	Image       string `json:"image,omitempty" yaml:"image,omitempty"`
	DockerImage string `json:"dockerImage,omitempty" yaml:"dockerImage,omitempty"`

	// (required) Container path in which the task will save results. Files
	// written to this location will be persisted as a dataset upon task
	// completion.
	ResultPath string `json:"resultPath" yaml:"resultPath"`

	// (optional) Text description of the task.
	Description string `json:"desc" yaml:"description,omitempty"` // TODO: Rename to "description"

	// (optional) Entrypoint to pass to the task's container.
	Command []string `json:"command,omitempty" yaml:"command,omitempty"`

	// (optional) Command-line arguments to pass to the task's container.
	Arguments []string `json:"arguments,omitempty" yaml:"args,omitempty"`

	// (optional) Environment variables to pass into the task's container.
	Env map[string]string `json:"env" yaml:"env,omitempty"`

	// TODO: Replace both mount lists with TaskMount.

	// (optional) Data sources to mount as read-only in the task's container.
	// In the event that mounts overlap partially or in full, they will be
	// applied in order. Later mounts will overlay earlier ones (last wins).
	Mounts []DatasetMount `json:"sources" yaml:"datasetMounts,omitempty"` // TODO: Rename to "mounts"

	// (optional) Task resource requirements for scheduling.
	Requirements TaskRequirements `json:"requirements" yaml:"requirements,omitempty"`
}

// TaskRequirements describes the runtime hardware requirements for a task.
type TaskRequirements struct {
	// (optional) Minimum required memory, in bytes.
	Memory int64 `json:"memory" yaml:"-"`

	// (optional) Minimum required memory, as a string which includes unit suffix.
	// Examples: "2g", "256m"
	MemoryHuman string `json:"-" yaml:"memory,omitempty"`

	// (optional) Minimum CPUs to allocate in millicpus (1 CPU = 1000 millicpus).
	MilliCPU int `json:"cpu" yaml:"-"`

	// (optional) Minimum CPUs to allocate as floating point.
	// CPU requirements are rounded to one thousandth of a CPU, i.e. 0.001
	CPU float64 `json:"-" yaml:"cpu,omitempty"`

	// (optional) GPUs required in increments of one full core.
	GPUCount int `json:"gpuCount" yaml:"gpuCount,omitempty"`
}

// DatasetMount describes a read-only data source for a task.
type DatasetMount struct {
	// (required) Name or Unique ID of a dataset to mount.
	Dataset string `json:"dataset" yaml:"datasetId"` // YAML name is locked for v1 spec.

	// (optional) Path within the dataset to mount for this experiment container.
	SubPath string `json:"subPath,omitempty" yaml:"subPath,omitempty"`

	// (required) Path within a task container to which file(s) will be mounted.
	ContainerPath string `json:"containerPath" yaml:"containerPath"`
}

// TaskPatchSpec describes a patch to apply to a task's editable fields.
type TaskPatchSpec struct {
	// (optional) Description to assign to the task or empty string to delete an
	// existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Whether the task should be canceled. Ignored if false.
	Cancel bool `json:"cancel,omitempty"`
}

type TaskEvents struct {
	Task   string      `json:"task"`
	Events []TaskEvent `json:"events"`
}

type TaskEvent struct {
	Status  TaskStatus `json:"status"`
	Message string     `json:"message,omitempty"`
	Time    time.Time  `json:"time"`
}
