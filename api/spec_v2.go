package api

import "github.com/allenai/bytefmt"

// ExperimentSpecV2 describes a collection of processes to run.
type ExperimentSpecV2 struct {
	// (required) Version must be 'v2-alpha'
	Version string `json:"version" yaml:"version"`

	// (optional) Description provides a long-form explanation for an experiment.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// (required) Tasks define what to run.
	Tasks []TaskSpecV2 `json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

// TaskSpecV2 describes a single job, or process, to run.
type TaskSpecV2 struct {
	// (optional) Name is used to refer to this task. It must be unique among
	// all tasks within the Spec.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// (required) Image is the name or ID of an image to run.
	Image ImageSource `json:"image" yaml:"image"`

	// (optional) Command is the full shell command to run as a list of separate
	// arguments. If omitted, the image's default command is used, for example
	// Docker's ENTRYPOINT directive. If set, default commands such as Docker's
	// ENTRYPOINT and CMD directives are ignored.
	//
	// Example: ["python", "-u", "main.py"]
	Command []string `json:"command,omitempty" yaml:"command,omitempty,flow"`

	// (optional) Arguments are appended to the Command and replace default
	// arguments such as Docker's CMD directive. If Command is omitted arguments
	// are appended to the default command, Docker's ENTRYPOINT directive.
	//
	// Example: If Command is ["python"], specifying arguments ["-u", "main.py"]
	// will run the command "python -u main.py".
	Arguments []string `json:"arguments,omitempty" yaml:"arguments,omitempty,flow"`

	// (optional) EnvVars are passed into the task as environment variables.
	EnvVars []EnvironmentVariable `json:"envVars,omitempty" yaml:"envVars,omitempty"`

	// (optional) Datasets are external data sources to mount into the task.
	Datasets []DataMount `json:"datasets,omitempty" yaml:"datasets,omitempty"`

	// (required) Result describes where the task will write output.
	Result ResultSpec `json:"result" yaml:"result"`

	// (optional) Resources define external hardware requirements for this task.
	// TODO: Consider whether to move this into the context.
	Resources *TaskResources `json:"resources,omitempty" yaml:"resources,omitempty"`

	// (required) Context describes how and where this task should run.
	//
	// Because contexts depend on external configuration, a given context may be
	// invalid or unavailable on subsequent runs.
	Context Context `json:"context" yaml:"context"`

	// (deprecated) Description is a long-form explanation of the task.
	Description string `json:"-" yaml:"-"`
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

// EnvironmentVariable describes the name and source of an environment variable.
// Exactly one source (value or secret) must be specified.
type EnvironmentVariable struct {
	// (required) Name of the environment variable. Case sensitive.
	Name string `json:"name" yaml:"name"`

	// (optional) Source the environment variable from a literal value.
	Value *string `json:"value,omitempty" yaml:"value,omitempty"`

	// (optional) Source the environment variable from a secret.
	// The secret must be present in the experiment's workspace.
	Secret string `json:"secret,omitempty" yaml:"secret,omitempty"`
}

// DataMount describes how a dataset is mounted into a task or environment.
type DataMount struct {
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
	// Source data from a Beaker dataset by name or ID.
	Beaker string `json:"beaker,omitempty" yaml:"beaker,omitempty"`

	// Source data from a host path. Support depends on the executing environment.
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`

	// Source data from of a previous task by name.
	Result string `json:"result,omitempty" yaml:"result,omitempty"`

	// Source data from a cloud service provider like S3/GS or HTTP
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// Source data from a secret. The secret is mounted as a file.
	// The secret must be in the same workspace as the experiment using it.
	Secret string `json:"secret,omitempty" yaml:"secret,omitempty"`
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
	CPUCount float64 `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`

	// (optional) GPUCount sets a mimimum number of GPU cores and must be a non-negative integer.
	GPUCount int `json:"gpuCount,omitempty" yaml:"gpuCount,omitempty"`

	// (optional) Memory sets a limit for CPU memory, which may be a raw number
	// of bytes or a formatted string with a number followed by a unit suffix.
	//
	// Examples: "2.5 GiB", 2684354560
	Memory *bytefmt.Size `json:"memory,omitempty" yaml:"memory,omitempty"`

	// (deprecated) Use Memory.
	// TODO: Remove this before resolving allenai/beaker-service#1106
	MemoryBytes int64 `json:"memoryBytes,omitempty" yaml:"memoryBytes,omitempty"`
}

// A Context describes how and where to run tasks.
type Context struct {
	// (required) Name or ID of a cluster on which the task should run.
	Cluster string `json:"cluster" yaml:"cluster"`

	// (optional) Priority describes the urgency with which a task will run.
	//
	// Values may be "low", "normal", or "high". If omitted, defaults to normal.
	// Priority may also be elevated to "urgent" through UI.
	Priority Priority `json:"priority,omitempty" yaml:"priority,omitempty"`
}

// Priority determines when a task will run relative to other tasks. Tasks of
// highest priority are exhausted before lower-priority tasks are considered for
// scheduling.
type Priority string

const (
	// UrgentPriority tasks are run as soon as possible. Because this priority
	// can be disruptive, it may only be set after a task's creation.
	UrgentPriority Priority = "urgent"

	// HighPriority tasks run before all non-urgent tasks. This is the highest
	// priority settable during task creation.
	HighPriority Priority = "high"

	// NormalPriority is the default priority for tasks, equivalent to unset.
	NormalPriority Priority = "normal"

	// LowPriority tasks run last and may be deferred for long periods of time.
	LowPriority Priority = "low"
)
