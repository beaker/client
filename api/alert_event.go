package api

import "time"

type AlertEventType string

const (
	AlertEventTypeTriggered AlertEventType = "triggered"
	AlertEventTypeCleared   AlertEventType = "cleared"
	AlertEventTypeEnabled   AlertEventType = "enabled"
	AlertEventTypeDisabled  AlertEventType = "disabled"
)

type AlertEvent struct {
	Event AlertEventType `json:"eventType"`
	Time  time.Time      `json:"time"`
}
