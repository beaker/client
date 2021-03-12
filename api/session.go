package api

// Session is an interactive Beaker session.
type Session struct {
	// Identity
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Cluster string   `json:"cluster"`
	Account string   `json:"account"` // Account of the session's cluster.
	Author  Identity `json:"author"`

	// Node that the session is assigned to.
	Node string `json:"node"`

	// Limits describes resources assigned to this session.
	Limits *TaskResources `json:"limits"`

	// State describes session status and progression.
	State ExecutionState `json:"state"`
}

// SessionSpec defines a session.
type SessionSpec struct {
	// (required) ID of the node that the session is running on.
	Node string `json:"node" yaml:"node"`

	// (optional) Name for the session.
	Name string `json:"name" yaml:"name"`
}

// SessionPatch updates a session.
type SessionPatch struct {
	// Limits updates the resources assigned to a session.
	Limits *TaskResources `json:"limits"`

	// State updates the session status and progression.
	State *ExecutionState `json:"state"`
}
