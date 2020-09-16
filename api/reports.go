package api

import (
	"time"

	"github.com/shopspring/decimal"
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

	// Whether or not the cluster is on-premise. Included if the data is grouped by on-premise.
	OnPrem *bool `json:"onPrem,omitempty"`

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

	// Whether or not the cluster is on-premise. Included if the data is grouped by on-premise.
	OnPrem *bool `json:"onPrem,omitempty"`

	Totals    UsageInterval   `json:"totals"`
	Intervals []UsageInterval `json:"intervals"`
}

// UsageReport:
// Deprecated: UsageReport - Reporting in general needs an improved API
type UsageReport struct {
	Start        time.Time     `json:"start"`
	End          time.Time     `json:"end"`
	Currency     string        `json:"currency"`
	EntityType   string        `json:"entityType"`
	Interval     string        `json:"interval"`
	IntervalKeys []time.Time   `json:"intervalKeys"`
	Items        []EntityUsage `json:"items"`
}

// EntityUsage:
// Deprecated: Reporting in general needs an improved API
type EntityUsage struct {
	Entity      Identity                          `json:"entity"`
	ReportGroup string                            `json:"reportGroup"`
	Totals      LegacyUsageInterval               `json:"totals"`
	Intervals   map[time.Time]LegacyUsageInterval `json:"intervals"`
}

// LegacyUsageInterval:
// Deprecated: Reporting in general needs an improved API
type LegacyUsageInterval struct {
	ExperimentCount int             `json:"experimentCount"`
	Cost            decimal.Decimal `json:"cost"`         // Total cost, i.e. UsageCost + OverheadCost
	UsageCost       decimal.Decimal `json:"usageCost"`    // Distinguish actual usage related cost from total cost
	OverheadCost    decimal.Decimal `json:"overheadCost"` // Distinguish overhead
	Duration        int64           `json:"duration"`
	GPUSeconds      int64           `json:"gpuSeconds"`
}
