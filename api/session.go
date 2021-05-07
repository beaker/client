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
	Node string `json:"node,omitempty"`

	// State describes session status and progression.
	State ExecutionState `json:"state"`

	// Resources requested by this session.
	Requests *ResourceRequest `json:"requests,omitempty"`

	// Limits assigned to this session.
	Limits *ResourceLimits `json:"limits,omitempty"`
}

// SessionSpec defines a session.
type SessionSpec struct {
	// (required) ID of the node that the session is running on.
	Node string `json:"node"`

	// (optional) Name for the session.
	Name string `json:"name"`

	// (optional) Resources requested by this session.
	Requests *ResourceRequest `json:"requests,omitempty"`
}

// SessionPatch updates a session.
type SessionPatch struct {
	// (optional) State updates the session status and progression.
	State *ExecStatusUpdate `json:"state"`

	// (optional) Limits updates the resources assigned to a session.
	Limits *ResourceLimits `json:"limits"`
}
