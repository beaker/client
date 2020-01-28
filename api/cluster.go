package api

import (
	"time"

	"github.com/shopspring/decimal"
)

// A Cluster is a homogenous collection of compute nodes. Nodes may be
// virtual machines or physical hardware, depending on the hosting environment.
type Cluster struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`

	Created    time.Time  `json:"created"`
	Terminated *time.Time `json:"terminated,omitempty"`

	// Capacity is the maximum number of nodes a cluster can contain at one time.
	Capacity int `json:"capacity"`

	// NodeCost describes the cost per node in units of USD-per-hour.
	NodeCost *decimal.Decimal `json:"nodeCost,omitempty"`

	// Requested and actual configuration
	Preemptible bool          `json:"preemptible"`
	Protected   bool          `json:"protected"`
	Status      ClusterStatus `json:"status"`
	NodeSpec    NodeSpec      `json:"nodeSpec"`
	NodeShape   *NodeSpec     `json:"nodeShape,omitempty"`
}

// ClusterPatch allows a client to update aspects of a Cluster.
type ClusterPatch struct {
	// Capacity changes the maximum number of nodes a cluster can contain at one time.
	Capacity *int `json:"capacity,omitempty"`

	// Valid permanently sets validity for the cluster and should be accompanied
	// by an node spec in the same request. If set to true, the cluster is ready
	// for use. Otherwise, it's considered failed.
	//
	// This value is internal; behavior is undefined if set by external clients.
	Valid *bool `json:"valid,omitempty"`

	// NodeShape details the shape of nodes created during cluster creation.
	//
	// This value is internal; behavior is undefined if set by external clients.
	NodeShape *NodeSpec `json:"nodeShape,omitempty"`

	// NodeCost sets the estimated cost of each node within the cluster in USD-per-hour.
	NodeCost *decimal.Decimal `json:"nodeCost,omitempty"`
}

// A ClusterPage contains a partial list of clusters.
type ClusterPage struct {
	Data       []Cluster `json:"data"`
	NextCursor string    `json:"nextCursor"`
}

// ClusterSpec provides options to configure a new cluster.
type ClusterSpec struct {
	Name     string `json:"name,omitempty"`
	Capacity int    `json:"capacity"`

	// Preemptible declares whether the cluster should include lower cost
	// preemptible nodes, with the tradeoff that workloads are more likely to be
	// interrupted.
	Preemptible bool `json:"preemptible,omitempty"`

	// Protected marks the cluster as a protected resource, ensuring only
	// administrators can modify/delete it. Only administrators can create them.
	Protected bool `json:"protected,omitempty"`

	// Spec describes characteristics of each node within the cluster.
	// Default values will be set by internal policy.
	Spec NodeSpec `json:"spec"`
}

// ClusterStatus describes where a cluster is in its lifecycle.
type ClusterStatus string

const (
	// ClusterPending indicates a cluster is in the process of being created.
	ClusterPending ClusterStatus = "pending"

	// ClusterActive indicates a cluster is online and available to schedule tasks.
	ClusterActive ClusterStatus = "active"

	// ClusterTerminated indicates a cluster has expired or been explicitly stopped.
	ClusterTerminated ClusterStatus = "terminated"

	// ClusterFailed indicates creation of a cluster could not be completed.
	ClusterFailed ClusterStatus = "failed"
)

// NodeSpec provides options to configure compute nodes.
type NodeSpec struct {
	CPUCount int    `json:"cpuCount,omitempty"`
	GPUCount int    `json:"gpuCount,omitempty"`
	GPUType  string `json:"gpuType,omitempty"`
	Memory   string `json:"memory,omitempty"`
}

// A Node is a single machine within a cluster
type Node struct {
	ID         string     `json:"id"`
	Hostname   string     `json:"hostname"`
	Created    time.Time  `json:"created"`
	Terminated *time.Time `json:"terminated"`
}

// NodePage contains a partial list of nodes.
type NodePage struct {
	Data []Node `json:"data"`
}

// CreateNodeSpec allows a requestor to describe a node on creation.
type CreateNodeSpec struct {
	Hostname string `json:"hostname"`
}

// ScheduledTasks is an ordered collection of scheduled tasks.
type ScheduledTasks struct {
	Data []ScheduledTask `json:"data"`
}

// Executions is an ordered collection of task executions.
type Executions struct {
	Data []Execution `json:"data"`
}

// Execution represents an attempt to run a task. A task may have many executions.
type Execution struct {
	ID   string `json:"id"`
	Task string `json:"task"`

	Image   ImageSource           `json:"image"`
	Result  ResultTarget          `json:"result"`
	Sources map[string]DataSource `json:"sources"` // Keyed by container path

	// Node is set when a task has been assigned to a node.
	Node string `json:"node,omitempty"`

	// State describes execution status and progression.
	State ExecutionState `json:"state"`
}

// ResultTarget describes a target to which results will be written.
type ResultTarget struct {
	// Name or ID of a Beaker dataset.
	Beaker string `json:"beaker,omitempty"`
}

// ExecutionState details an execution's status.
type ExecutionState struct {
	Created   time.Time  `json:"created"`
	Scheduled *time.Time `json:"scheduled,omitempty"`
	Started   *time.Time `json:"started,omitempty"`
	Ended     *time.Time `json:"ended,omitempty"`
	Finalized *time.Time `json:"finalized,omitempty"`

	// ExitCode is an integer process exit code, if the process exited normally.
	ExitCode *int `json:"exitCode,omitempty"`

	// Message describes additional state-related context.
	Message string `json:"message,omitempty"`

	// TaskCanceled indicates whether and when an execution's task was canceled.
	TaskCanceled *time.Time `json:"taskCanceled,omitempty"`
}

// ScheduledTask summarizes relations of executions, or tasks which have been scheduled to run.
// TODO: Remove this and replace with Execution.
type ScheduledTask struct {
	Task string `json:"task"`

	// Execution uniquely identifies one attempt to execute a task.
	ExecutionID string `json:"id"` // Back-compat shim so we can update executors.
	Execution   string `json:"execution"`

	// Node is set when a task has been assigned to a node.
	Node string `json:"node,omitempty"`
}
