package api

type TFEventNameList struct {
	Events []string `json:"events"`
}

type TFEventSeriesList struct {
	Series []TFEventSeries `json:"series"`
}

type TFEventSeries struct {
	Experiment string `json:"experiment"`
	// File from which the events were extracted relative to the task's result directory.
	Path string `json:"path"`
	Task string `json:"task"`

	// All arrays are guaranteed to have the same length.
	Steps  []int64   `json:"steps"`
	Times  []int64   `json:"times"` // Seconds since the Unix epoch.
	Values []float32 `json:"values"`
}
