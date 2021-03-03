package api

// Session is an interactive Beaker session.
type Session struct {
	ID string `json:"id"`

	// Limits describes resources assigned to this session.
	Limits TaskResources `json:"limits"`

	// State describes session status and progression.
	State ExecutionState `json:"state"`
}
