package api

import (
	"time"
)

// Interval specifies the granularity of each interval in a report.
type Interval string

// Reports can be broken down by these intervals
const (
	IntervalHour    Interval = "hour"
	IntervalDay     Interval = "day"
	IntervalWeek    Interval = "week"
	IntervalMonth   Interval = "month"
	IntervalQuarter Interval = "quarter"
	IntervalYear    Interval = "year"
)

// NodeUsageReport contains one series for each combination of values in the group by.
type NodeUsageReport struct {
	Totals UsageInterval     `json:"totals"`
	Series []NodeUsageSeries `json:"series"`
}

// UsageInterval reports the value of a NodeMetric during a single time interval.
type UsageInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Value float64   `json:"value"`
}

// NodeUsageSeries reports the value of a NodeMetric over time.
type NodeUsageSeries struct {
	// Node ID is included if the data is grouped by node.
	Node string `json:"node,omitempty"`

	// Cluster name is included if the data is grouped by cluster.
	Cluster string `json:"cluster,omitempty"`

	// Whether or not the cluster is preemptible. Included if the data is grouped by preemptible.
	Preemptible *bool `json:"preemptible,omitempty"`

	// Whether the node's cluster scales automatically. Included if the data is grouped by autoscale.
	Autoscale *bool `json:"autoscale,omitempty"`

	// Number of GPUs the node has. Included if the data is grouped by GPU count.
	GPUCount *int `json:"gpuCount,omitempty"`

	// Type of GPU that the node has. Included if the data is grouped by GPU type.
	GPUType *string `json:"gpuType,omitempty"`

	Totals    UsageInterval   `json:"totals"`
	Intervals []UsageInterval `json:"intervals"`
}

// TaskUsageReport contains one series for each combination of values in the group by.
type TaskUsageReport struct {
	Totals UsageInterval     `json:"totals"`
	Series []TaskUsageSeries `json:"series"`
}

// TaskUsageSeries reports the value of a TaskMetric over time.
type TaskUsageSeries struct {
	// Task ID is included if the data is grouped by node.
	Task string `json:"task,omitempty"`

	// Experiment ID is included if the data is grouped by experiment.
	Experiment string `json:"experiment,omitempty"`

	// Workspace ID is included if the data is grouped by workspace.
	Workspace string `json:"workspace,omitempty"`

	// Node ID is included if the data is grouped by node.
	Node string `json:"node,omitempty"`

	// Cluster name is included if the data is grouped by cluster.
	Cluster string `json:"cluster,omitempty"`

	// Author name is included if the data is grouped by author.
	Author string `json:"author,omitempty"`

	// Owner name is included if the data is grouped by owner.
	Owner string `json:"owner,omitempty"`

	// Team name is included if the data is grouped by team.
	Team string `json:"team,omitempty"`

	// Whether the cluster running the task autoscales. Included if the data is grouped by autoscale.
	Autoscale *bool `json:"autoscale,omitempty"`

	// Number of GPUs claimed by the task. Included if the data is grouped by GPU count.
	GPUCount *int `json:"gpuCount,omitempty"`

	// Type of GPU that the task's node has. Included if the data is grouped by GPU type.
	GPUType *string `json:"gpuType,omitempty"`

	Totals    UsageInterval   `json:"totals"`
	Intervals []UsageInterval `json:"intervals"`
}
