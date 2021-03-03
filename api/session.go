package api

// Session is an interactive Beaker session.
type Session struct {
	ID string `json:"id"`

	// Node that the session is assigned to.
	Node string `json:"node"`

	// Limits describes resources assigned to this session.
	Limits TaskResources `json:"limits"`

	// State describes session status and progression.
	State ExecutionState `json:"state"`
}

// SessionSpec defines a session.
type SessionSpec struct {
	// (required) ID of the node that the session is running on.
	Node string `json:"node" yaml:"node"`
}

// SessionPatch updates a session.
type SessionPatch struct {
	// Limits updates the resources assigned to a session.
	Limits *TaskResources `json:"limits"`

	// State updates the session status and progression.
	State *ExecutionState `json:"state"`
}
