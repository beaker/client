package api

import "github.com/allenai/bytefmt"

// Session is an interactive Beaker session.
type Session struct {
	// Identity
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Cluster string   `json:"cluster"`
	Account string   `json:"account"` // Account of the session's cluster.
	Author  Identity `json:"author"`

	// Node that the session is assigned to.
	Node string `json:"node,omitempty"`

	// State describes session status and progression.
	State ExecutionState `json:"state"`

	// Resources requested by this session.
	Requests *TaskResources `json:"requests,omitempty"`

	// Limits assigned to this session.
	Limits *SessionResources `json:"limits,omitempty"`
}

// SessionSpec defines a session.
type SessionSpec struct {
	// (required) ID of the node that the session is running on.
	Node string `json:"node"`

	// (optional) Name for the session.
	Name string `json:"name"`

	// (optional) Resources requested by this session.
	Requests *TaskResources `json:"requests,omitempty"`
}

// SessionPatch updates a session.
type SessionPatch struct {
	// (optional) State updates the session status and progression.
	State *ExecStatusUpdate `json:"state"`

	// (optional) Limits updates the resources assigned to a session.
	Limits *SessionResources `json:"limits"`
}

// SessionResources describe external requirements which must be available for a session to run.
type SessionResources struct {
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
