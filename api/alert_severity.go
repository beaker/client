package api

// AlertSeverity represents the severity level for an alert
type AlertSeverity string

const (
	AlertSeverityLow AlertSeverity = "low"

	AlertSeverityMed AlertSeverity = "med"

	AlertSeverityHigh AlertSeverity = "high"
)
