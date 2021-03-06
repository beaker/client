package api

import (
	"time"

	"github.com/allenai/bytefmt"
	"github.com/shopspring/decimal"
)

// A Cluster is a homogenous collection of compute nodes. Nodes may be
// virtual machines or physical hardware, depending on the hosting environment.
type Cluster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	Created   time.Time `json:"created"`
	Protected bool      `json:"protected"`

	// Everything after this point is autoscale policy.
	// TODO allenai/beaker-service#1203: Separate autoscale from Cluster.

	Autoscale bool       `json:"autoscale"`
	Validated *time.Time `json:"validated,omitempty"`

	StatusMessage string `json:"statusMessage,omitempty"`

	// Capacity is the maximum number of nodes a cluster can contain at one time.
	Capacity int `json:"capacity"`

	// NodeCost describes the cost per node in units of USD-per-hour.
	NodeCost *decimal.Decimal `json:"nodeCost,omitempty"`

	// Requested and actual configuration
	Preemptible bool           `json:"preemptible"`
	Status      ClusterStatus  `json:"status"`
	NodeSpec    NodeResources  `json:"nodeSpec"`
	NodeShape   *NodeResources `json:"nodeShape,omitempty"`
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

	// StatusMessage provides optional additional information regarding the status,
	// e.g. why validation failed
	//
	// This value is internal; behavior is undefined if set by external clients.
	StatusMessage *string `json:"statusMessage,omitempty"`

	// NodeShape details the shape of nodes created during cluster creation.
	//
	// This value is internal; behavior is undefined if set by external clients.
	NodeShape *NodeResources `json:"nodeShape,omitempty"`

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
	// If omitted, the cluster will not autoscale.
	Spec *NodeResources `json:"spec"`
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

// NodeResources describe a node's available resources.
type NodeResources struct {
	CPUCount float64       `json:"cpuCount,omitempty"`
	GPUCount int           `json:"gpuCount,omitempty"`
	GPUType  string        `json:"gpuType,omitempty"`
	Memory   *bytefmt.Size `json:"memory,omitempty"`
}

// A Node is a single machine within a cluster
type Node struct {
	ID       string         `json:"id"`
	Hostname string         `json:"hostname"`
	Created  time.Time      `json:"created"`
	Expiry   *time.Time     `json:"expiry,omitempty"`
	Cordoned *time.Time     `json:"cordoned,omitempty"`
	Limits   *NodeResources `json:"limits,omitempty"`
}

// NodePage contains a partial list of nodes.
type NodePage struct {
	Data []Node `json:"data"`
}

// NodeSpec allows a requestor to describe a node on creation.
type NodeSpec struct {
	Hostname string         `json:"hostname"`
	Limits   *NodeResources `json:"limits"`
}

// NodePatchSpec allows a requestor to update a node.
// Expired nodes cannot be patched.
type NodePatchSpec struct {
	// (optional) Set to true to cordon the node or false to uncordon it.
	Cordoned *bool `json:"cordoned,omitempty"`

	// (optional) Formatted as a Golang duration e.g. "1h".
	// If included, the node's expiry is set to TTL after the current time.
	// TTL cannot exceed one week.
	TTL *string `json:"ttl,omitempty"`

	// (optional) If included, the node's hostname is updated.
	Hostname *string `json:"hostname,omitempty"`

	// (optional) If included, the node's limits are updated.
	Limits *NodeResources `json:"limits,omitempty"`
}
