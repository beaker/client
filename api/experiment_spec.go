package api

// ExperimentSpecV2 describes a collection of processes to run.
type ExperimentSpecV2 struct {
	// (required) Version must be 'v2'
	Version string `json:"version" yaml:"version"`

	// (required) Tasks to run.
	Tasks []TaskSpec `json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

// TaskSpecV2 describes a single job, or process, to run.
type TaskSpecV2 struct {
	// (optional) Name is used to refer to this task. It must be unique among
	// all tasks within the Spec.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// (required) Image is the name or ID of an image to run.
	Image ImageSource `json:"image" yaml:"image"`

	// (required) Command is the full shell command to run as a list of separate
	// arguments. Default comamnds such as Docker's ENTRYPOINT are ignored.
	// Example: ["python", "-u", "main.py"]
	Command []string `json:"command" yaml:"command"`

	// (optional) EnvVars are passed into the task as environment variables.
	EnvVars map[string]string `json:"envVars,omitempty" yaml:"envVars,omitempty"`

	// (optional) Datasets are external data sources to mount into the task.
	Datasets []Mount `json:"datasets,omitempty" yaml:"datasets,omitempty"`

	// (required) Result describes where the task will write output.
	Result ResultSpec `json:"result" yaml:"result"`

	// (optional) Resources define external requirements for task execution.
	Resources *TaskResources `json:"resources,omitempty" yaml:"resources,omitempty"`
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

// DataSource describes all supported data sources by type. Exactly one must be defined.
type DataSource struct {
	// Name or ID of a Beaker dataset.
	Beaker string `json:"beaker,omitempty" yaml:"beaker,omitempty"`

	// Mount data from the executing host. Support depends on the executing environment.
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`

	// Name of a previous task from which the result will be mounted.
	Result string `json:"result,omitempty" yaml:"result,omitempty"`
}

// ResultSpec describes how to store the output of a task.
type ResultSpec struct {
	// (required) Path is a file or directory where the task will write output.
	Path string `json:"path" yaml:"path"`
}

// TaskResources describe external requirements which must be available for a task to run.
type TaskResources struct {
	// (optional) CPUCount sets a minimum number of logical CPU cores and may be fractional.
	//
	// Examples: 4, 0.5
	CPUCount float64 `json:"cpuCount,omitempty" yaml:"gpuCount,omitempty"`

	// (optional) GPUCount sets a mimimum number of GPU cores and must be a non-negative integer.
	GPUCount int `json:"gpuCount,omitempty" yaml:"gpuCount,omitempty"`

	// (optional) Memory sets a limit for CPU memory as a number with unit suffix.
	//
	// Examples: 2.5GiB, 10240m
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}
