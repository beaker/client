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

	// State of this task and its execution(s).
	Created    time.Time       `json:"created"`
	Canceled   *time.Time      `json:"canceled,omitempty"`
	LastState  *ExecutionState `json:"lastState,omitempty"`
	Executions []string        `json:"executions,omitempty"`

	// Creation parameters
	Spec        TaskSpecV1      `json:"spec"`
	SpecV1      *TaskSpecV1     `json:"specV1,omitempty"`
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

// TaskPatchSpec describes a patch to apply to a task's editable fields.
type TaskPatchSpec struct {
	// (optional) Description to assign to the task or empty string to delete an
	// existing description.
	Description *string `json:"description,omitempty"`

	// (optional) Whether the task should be canceled. Ignored if false.
	Cancel bool `json:"cancel,omitempty"`
}
