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
	Expiration *time.Time `json:"expiration,omitempty"`
	Terminated *time.Time `json:"terminated,omitempty"`

	// Capacity is the maximum number of nodes a cluster can contain at one time.
	Capacity int `json:"capacity"`

	// NodeCost describes the cost per node in units of USD-per-hour.
	NodeCost *decimal.Decimal `json:"nodeCost,omitempty"`

	// Requested and actual configuration
	Status    ClusterStatus `json:"status"`
	NodeSpec  NodeSpec      `json:"nodeSpec"`
	NodeShape *NodeSpec     `json:"nodeShape,omitempty"`
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
	Data []Cluster `json:"data"`
}

// ClusterSpec provides options to configure a new cluster.
type ClusterSpec struct {
	Name     string `json:"name,omitempty"`
	Capacity int    `json:"capacity"`

	// Preemptible declares whether the cluster should include lower cost
	// preemptible nodes, with the tradeoff that workloads are more likely to be
	// interrupted.
	Preemptible bool `json:"preemptible,omitempty"`

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

// Deprecated. Use clusters/nodes instead.
type Machine struct {
	ID          string `json:"id"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	NodeLabel   string `json:"nodeLabel"`
	GPUCount    int    `json:"gpuCount,omitempty"`
	GPUType     string `json:"gpuType,omitempty"`
	GPULabel    string `json:"gpuLabel,omitempty"`
	Preemptible bool   `json:"preemptible"`
	Cost        int64  `json:"cost"`
	IsActive    bool   `json:"isActive"`
}

// NodeSpec provides options to configure compute nodes.
type NodeSpec struct {
	CPUCount int    `json:"cpuCount,omitempty"`
	GPUCount int    `json:"gpuCount,omitempty"`
	GPUType  string `json:"gpuType,omitempty"`
	Memory   string `json:"memory,omitempty"`
}

// NodeStatus describes the availability of a node.
type NodeStatus string

const (
	// Starting nodes are in the process of booting.
	Starting NodeStatus = "starting"

	// Running nodes are online.
	Running NodeStatus = "running"

	// Stopping nodes are shutting down.
	Stopping NodeStatus = "stopping"

	// Stopped nodes have been shut down, but not terminated.
	Stopped NodeStatus = "stopped"

	// Terminated nodes have been permanently stopped (deleted).
	Terminated NodeStatus = "terminated"
)

// A Node is a single machine within a cluster
type Node struct {
	ID       string     `json:"id"`
	Hostname string     `json:"hostname"`
	Created  time.Time  `json:"created"`
	Status   NodeStatus `json:"status"`
}

// NodePage contains a partial list of nodes.
type NodePage struct {
	Data []Node `json:"data"`
}

// CreateNodeSpec allows a requestor to describe a node on creation.
type CreateNodeSpec struct {
	Hostname string `json:"hostname"`
}

// NodeStatusSpec describes a status change for a node.
type NodeStatusSpec struct {
	Status NodeStatus `json:"status"`
}

// ScheduledTasks is an ordered collection of task executions.
type ScheduledTasks struct {
	Data []ScheduledTask `json:"data"`
}

// ScheduledTask summarizes relations of tasks which have been scheduled.
type ScheduledTask struct {
	TaskID      string `json:"taskId"`
	ExecutionID string `json:"executionId"`

	// NodeID is set when a task has been scheduled on an node.
	NodeID string `json:"nodeId,omitempty"`
}
