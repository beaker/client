package api

// ExperimentSpec describes a set of tasks with optional dependencies.
// This set represents a (potentially disconnected) directed acyclic graph.
type ExperimentSpecV1 struct {
	// (optional) Version must be 'v1' or left unset.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// (required) Workspace where this experiment and its results should be placed.
	Workspace string `json:"workspace,omitempty" yaml:"workspace,omitempty"`

	// (optional) Text description of the experiment.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// (required) Tasks to create. Tasks may be defined in any order, though all
	// dependencies must be internally resolvable within the experiment.
	Tasks []ExperimentTaskSpec `json:"tasks" yaml:"tasks"`

	// (optional) A token representing the user to which the object should be attributed.
	// If omitted attribution will be given to the user issuing the request.
	AuthorToken string `json:"authorToken,omitempty" yaml:"authorToken,omitempty"`
}

// ExperimentTaskSpec describes a task spec with optional dependencies on other
// tasks within an experiment. Tasks refer to each other by the Name field.
type ExperimentTaskSpec struct {
	// (optional) Name of the task node, which need only be defined if
	// dependencies reference it.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// (required) Specification describing the task to run.
	Spec TaskSpecV1 `json:"spec" yaml:"spec,omitempty"`

	// (optional) Tasks on which this task depends. Mounts will be applied, in
	// the order defined here, after existing mounts in the task spec.
	DependsOn []TaskDependency `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`

	// (optional) Name of a cluster on which the task should run.
	// Cluster affinity supercedes task requirements.
	Cluster string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

// TaskSpecV1 contains all information necessary to create a new task.
type TaskSpecV1 struct {
	// (required) Image containing the code to be run.
	Image       string `json:"image,omitempty" yaml:"image,omitempty"`
	DockerImage string `json:"dockerImage,omitempty" yaml:"dockerImage,omitempty"`

	// (required) Container path in which the task will save results. Files
	// written to this location will be persisted as a dataset upon task
	// completion.
	ResultPath string `json:"resultPath" yaml:"resultPath"`

	// (optional) Text description of the task.
	Description string `json:"desc" yaml:"description,omitempty"`

	// (optional) Entrypoint to pass to the task's container.
	Command []string `json:"command,omitempty" yaml:"command,omitempty,flow"`

	// (optional) Command-line arguments to pass to the task's container.
	Arguments []string `json:"arguments,omitempty" yaml:"args,omitempty,flow"`

	// (optional) Environment variables to pass into the task's container.
	Env map[string]string `json:"env" yaml:"env,omitempty"`

	// (optional) Data sources to mount as read-only in the task's container.
	// In the event that mounts overlap partially or in full, they will be
	// applied in order. Later mounts will overlay earlier ones (last wins).
	Mounts []DatasetMount `json:"sources" yaml:"datasetMounts,omitempty"`

	// (optional) Task resource requirements for scheduling.
	Requirements TaskRequirements `json:"requirements" yaml:"requirements,omitempty"`
}

// DatasetMount describes a read-only data source for a task.
type DatasetMount struct {
	// (required) Name or Unique ID of a dataset to mount.
	Dataset string `json:"dataset" yaml:"datasetId"` // YAML name is locked for v1 spec.

	// (optional) Path within the dataset to mount for this experiment container.
	SubPath string `json:"subPath,omitempty" yaml:"subPath,omitempty"`

	// (required) Path within a task container to which file(s) will be mounted.
	ContainerPath string `json:"containerPath" yaml:"containerPath"`
}

// TaskRequirements describes the runtime hardware requirements for a task.
type TaskRequirements struct {
	// (optional) Minimum required memory, in bytes.
	Memory int64 `json:"memory" yaml:"-"`

	// (optional) Minimum required memory, as a string which includes unit suffix.
	// Examples: "2g", "256m"
	MemoryHuman string `json:"-" yaml:"memory,omitempty"`

	// (optional) Minimum CPUs to allocate in millicpus (1 CPU = 1000 millicpus).
	MilliCPU int `json:"cpu" yaml:"-"`

	// (optional) Minimum CPUs to allocate as floating point.
	// CPU requirements are rounded to one thousandth of a CPU, i.e. 0.001
	CPU float64 `json:"-" yaml:"cpu,omitempty"`

	// (optional) GPUs required in increments of one full core.
	GPUCount int `json:"gpuCount" yaml:"gpuCount,omitempty"`
}
