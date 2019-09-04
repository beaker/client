package api

// CreateMetricAlertResponse is a service response returned when a new metric alert is created.
type CreateMetricAlertResponse struct {
	ID string `json:"id"`
}

// MetricAlert is a full description of a metric alert specification
type MetricAlert struct {
	// ID
	ID     string `json:"id"`
	TaskID string `json:"taskId"` // Task emitting the metric to alert on

	// Ownership
	Owner  Identity `json:"owner"`
	Author Identity `json:"author"`

	// Metrics
	Metric    string               `json:"metric"`
	Condition MetricAlertCondition `json:"condition"`
	Threshold float64              `json:"threshold"`
	Severity  AlertSeverity        `json:"severity"`

	// Currently enabled and triggered status
	Enabled   bool `json:"enabled"`
	Triggered bool `json:"triggered"`
}

// MetricAlertSpec is currently for internal use only
type MetricAlertSpec struct {
	TaskID    string               `json:"taskId" yaml:"taskID"`
	Metric    string               `json:"metric" yaml:"metric"`
	Condition MetricAlertCondition `json:"condition" yaml:"condition"`
	Threshold float64              `json:"threshold" yaml:"threshold"`
	Severity  AlertSeverity        `json:"severity" yaml:"severity"`
}

// MetricAlertPatchSpec describes a patch to apply to a dataset's editable fields.
// Only one field may be set in a single request.
type MetricAlertPatchSpec struct {
	// (optional) Whether the alert should be (re-)enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// (optional) Whether the alert should be triggered.
	Triggered *bool `json:"triggered,omitempty"`
}