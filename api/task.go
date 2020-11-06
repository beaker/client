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
	Created    time.Time   `json:"created"`
	Canceled   *time.Time  `json:"canceled,omitempty"`
	Executions []Execution `json:"executions,omitempty"`

	// Creation parameters
	Spec        TaskSpecV2      `json:"spec"`
	SpecV1      *TaskSpecV1     `json:"specV1,omitempty"`
	ResumedFrom ResumedFromSpec `json:"resumedFrom"`

	// Deprecated
	FullExecutions []Execution `json:"fullExecutions,omitempty"`
}

type ResumedFromSpec struct {
	TaskID       string `json:"taskId,omitempty"`
	ExperimentID string `json:"experimentId,omitempty"`
}
