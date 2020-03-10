package api

import (
	"time"

	"github.com/shopspring/decimal"
)

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
