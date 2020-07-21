package api

// ExecutionPatchSpec describes a patch to apply to a execution's editable fields.
type ExecutionPatchSpec struct {
	// (optional) Priority to assign to the execution.
	Priority *string `json:"priority,omitempty"`
}
