package api

// Execution contains the information necessary for an executor to run a task.
type Execution struct {
	ID     string `json:"id"`
	TaskID string `json:"taskId"`

	// Image is a fully-resolved image source.
	Image ImageSource `json:"image"`

	// Command is the full command to run as a list of separate arguments.
	Command []string `json:"command"`

	// EnvVars is a mapping of all environment variables passed into the task.
	EnvVars map[string]string `json:"envVars,omitempty"`

	// Datasets are external data sources to mount into the task.
	Datasets []Mount `json:"datasets,omitempty"`

	// Result describes where the task will write output.
	Result ResultMount `json:"result"`

	// Resources define external requirements for task execution.
	Resources *TaskResources `json:"resources,omitempty"`
}

// ImageSource describes all supported image sources by type. Exactly one must be defined.
type ImageSource struct {
	// Name or ID of a Beaker image.
	Beaker string `json:"beaker,omitempty" yaml:"beaker,omitempty"`

	// Reference (SHA or name) of a local or remote Docker image, including
	// registry. If the image is from a private registry, the building host must
	// be pre-configured to allow access.
	Docker string `json:"docker,omitempty" yaml:"docker,omitempty"`
}

// Mount describes how a dataset is mounted into a task or environment.
type Mount struct {
	// (required) Path within a container to mount the data source. Mount paths
	// must be absolute and may not overlap.
	//
	// As some environments use case-insensitive file systems, mount paths
	// differing only in capitilization are considered overlapping.
	MountPath string `json:"mountPath" yaml:"mountPath"`

	// (optional) Sub-path to a file or directory within a mounted data source.
	// Subpaths may be used to mount a portion of a dataset; files outside of
	// the mounted path are not transferred.
	//
	// Example: For a dataset containing a file "/path/to/file.csv", a sub-path
	// of "/path/to" will show up to a task as "<mount-path>/file.csv".
	SubPath string `json:"subPath,omitempty" yaml:"subPath,omitempty"`

	// (required) Source describes where to find data to mount.
	Source DataSource `json:"source" yaml:"source"`
}

// DataSource describes all supported data sources by type. Exactly one type
// must be defined.
type DataSource struct {
	// Name or ID of a Beaker dataset.
	Beaker string `json:"beaker,omitempty" yaml:"beaker,omitempty"`

	// Mount data from the executing host. Support depends on the environment.
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`

	// Name of a previous task from which the result will be mounted.
	Result string `json:"result,omitempty" yaml:"result,omitempty"`
}

// ResultMount describes how to mount a result dataset within a task.
type ResultMount struct {
	// Path is a file or directory where the task will write output.
	Path string `json:"path"`

	// ID of the dataset to write results to.
	Dataset string `json:"dataset"`
}

// TaskResources describe external requirements which must be available for a task to run.
type TaskResources struct {
	GPUCount int `json:"gpuCount,omitempty"`
}
