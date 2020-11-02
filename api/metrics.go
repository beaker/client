package api

type TFEventNameList struct {
	Events []string `json:"events"`
}

type TFEventSeriesList struct {
	Series []TFEventSeries `json:"series"`
}

type TFEventSeries struct {
	// File from which the events were extracted relative to the task's result directory.
	Path string `json:"path"`
	Task string `json:"task"`

	// All arrays are guaranteed to have the same length.
	Steps  []int64   `json:"steps"`
	Times  []int64   `json:"times"` // Seconds since the Unix epoch.
	Values []float32 `json:"values"`
}

type SystemMetricNameList struct {
	SystemMetrics []string `json:"systemMetrics"`
}

type SystemMetricSeriesList struct {
	Series []SystemMetricSeries `json:"series"`
}

type SystemMetricSeries struct {
	Task string `json:"task"`

	Tags map[string]string `json:"tags"`

	// All arrays are guaranteed to have the same length.
	Times  []int64   `json:"times"` // Seconds since the Unix epoch.
	Values []float32 `json:"values"`
}

type SystemMetricAggregateList struct {
	Metrics []SystemMetricAggregate `json:"metrics"`
}

type SystemMetricAggregate struct {
	Task       string                      `json:"task"`
	Aggregates map[AggregationType]float32 `json:"aggregates"`
}

type AggregationType string

const (
	AggregationTypeMax    AggregationType = "max"
	AggregationTypeMin    AggregationType = "min"
	AggregationTypeMean   AggregationType = "mean"
	AggregationTypeMedian AggregationType = "median"
	AggregationTypeMode   AggregationType = "mode"
	AggregationTypeCount  AggregationType = "count"
)
