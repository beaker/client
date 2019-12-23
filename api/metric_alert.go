package api

// CreateMetricAlertResponse is a service response returned when a new metric alert is created.
type CreateMetricAlertResponse struct {
	ID string `json:"id"`
}

// MetricAlertByTask is a description of an alert as it corresponds to the task of an experiment
type MetricAlertByTask struct {
	ExperimentName string        `json:"experimentName"`
	ExperimentID   string        `json:"experimentId"`
	TaskID         string        `json:"taskId"`
	Alerts         []MetricAlert `json:"alerts"`
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
	Metric    string        `json:"metric"`
	Condition Condition     `json:"condition"`
	Threshold float64       `json:"threshold"`
	Severity  AlertSeverity `json:"severity"`

	// Currently enabled and triggered status
	Enabled   bool `json:"enabled"`
	Triggered bool `json:"triggered"`
}

// MetricAlerts is a collection of metric alerts
type MetricAlerts struct {
	MetricAlerts []MetricAlert `json:"metricAlerts"`
}

// MetricAlertSpec is currently for internal use only
type MetricAlertSpec struct {
	// (optional) Organization on behalf of whom this resource is created. The
	// user issuing the request must be a member of the organization. If omitted,
	// the resource will be owned by the requestor.
	Organization string `json:"org,omitempty"`

	// (required) Task whose metric the alert is tracking.
	TaskID string `json:"taskId" yaml:"taskID"`

	// (required) Particular metric in the task the alert is tracking.
	Metric string `json:"metric" yaml:"metric"`

	// (required) Condition with the threshold on which the alert should be triggered.
	Condition Condition `json:"condition" yaml:"condition"`

	// (required) Threshold with the condition on which the alert should be triggered.
	Threshold float64 `json:"threshold" yaml:"threshold"`

	// (required) Severity of the alert.
	Severity AlertSeverity `json:"severity" yaml:"severity"`

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty"`
}

// MetricAlertPatchSpec describes a patch to apply to a dataset's editable fields.
// Only one field may be set in a single request.
type MetricAlertPatchSpec struct {
	// (optional) Whether the alert should be (re-)enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// (optional) Whether the alert should be triggered.
	Triggered *bool `json:"triggered,omitempty"`
}

// UserMetricAlertPage describes a page of metric alert data
type UserMetricAlertPage struct {
	Data []MetricAlertByTask `json:"data"`
}

// Condition represents the condition on which the alert should be triggerered for the metric
type Condition string

const (
	// ConditionLessThan means the alert should be triggered when the metric falls below the given threshold
	ConditionLessThan Condition = "less_than"

	// ConditionGreaterThan means the alert should be triggered when the metric exceeds the given threshold
	ConditionGreaterThan Condition = "greater_than"
)
