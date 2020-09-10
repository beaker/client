package api

import (
	"time"

	"github.com/shopspring/decimal"
)

// NodeUsageReport contains one series for each combination of values in the group by.
type NodeUsageReport struct {
	Totals NodeUsageInterval `json:"totals"`
	Series []NodeUsageSeries `json:"series"`
}

// NodeUsageSeries reports the value of a NodeMetric over time.
type NodeUsageSeries struct {
	// Node ID is included if the data is grouped by node.
	Node string `json:"node,omitempty"`

	// Cluster name is included if the data is grouped by cluster.
	Cluster string `json:"cluster,omitempty"`

	// Environment (on-prem or cloud) is included if the data is grouped by environment.
	Environment string `json:"environment,omitempty"`

	Totals    NodeUsageInterval   `json:"totals"`
	Intervals []NodeUsageInterval `json:"intervals"`
}

// NodeUsageInterval reports the value of a NodeMetric during a single time interval.
type NodeUsageInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Value float64   `json:"value"`
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
	Entity      Identity                    `json:"entity"`
	ReportGroup string                      `json:"reportGroup"`
	Totals      UsageInterval               `json:"totals"`
	Intervals   map[time.Time]UsageInterval `json:"intervals"`
}

// UsageInterval:
// Deprecated: Reporting in general needs an improved API
type UsageInterval struct {
	ExperimentCount int             `json:"experimentCount"`
	Cost            decimal.Decimal `json:"cost"`         // Total cost, i.e. UsageCost + OverheadCost
	UsageCost       decimal.Decimal `json:"usageCost"`    // Distinguish actual usage related cost from total cost
	OverheadCost    decimal.Decimal `json:"overheadCost"` // Distinguish overhead
	Duration        int64           `json:"duration"`
	GPUSeconds      int64           `json:"gpuSeconds"`
}
