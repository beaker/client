package api

// MetricAlertCondition represents the condition on which the alert should be triggerered for the metric
type MetricAlertCondition string

const (
	// MetricAlertConditionLessThan means the alert should be triggered when the metric falls below the given threshold
	MetricAlertConditionLessThan MetricAlertCondition = "less_than"

	// MetricAlertConditionGreaterThan means the alert should be triggered when the metric exceeds the given threshold
	MetricAlertConditionGreaterThan MetricAlertCondition = "greater_than"
)
