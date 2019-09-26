package api

import "time"

// A Cluster is a homogenous collection of compute instances. Instances may be
// virtual machines or physical hardware, depending on the hosting environment.
type Cluster struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Galaxy string `json:"galaxy"`

	Created    time.Time  `json:"created"`
	Expiration *time.Time `json:"expiration,omitempty"`
	Terminated *time.Time `json:"terminated,omitempty"`

	// Capacity is the maximum number of instances a cluster can contain at one time.
	Capacity int `json:"capacity"`

	// InstanceSpec describes per-instance configuration.
	InstanceSpec InstanceSpec `json:"instanceSpec"`
}

// A ClusterPage contains a partial list of clusters.
type ClusterPage struct {
	Data []Cluster `json:"data"`
}

// ClusterSpec provides options to configure a new cluster.
type ClusterSpec struct {
	Name     string `json:"name,omitempty"`
	Galaxy   string `json:"galaxy,omitempty"`
	Capacity int    `json:"capacity"`

	// Preemptible declares whether the cluster should include lower cost
	// preemptible instances, with the tradeoff that workloads are more likely
	// to be interrupted.
	Preemptible bool `json:"preemptible,omitempty"`

	// Spec describes characteristics of each instance within the cluster.
	// Default values will be set by internal policy.
	Spec InstanceSpec `json:"spec"`
}

// InstanceSpec provides options to configure compute instances.
type InstanceSpec struct {
	CPUCount int    `json:"cpuCount"`
	GPUCount int    `json:"gpuCount,omitempty"`
	GPUType  string `json:"gpuType,omitempty"`
	Memory   string `json:"memory"`
}

// InstanceSummary summarizes a instance's current status.
type InstanceSummary struct {
	ID     string         `json:"id"`
	Status InstanceStatus `json:"status"`

	// IDs of all tasks running or scheduled on this instance.
	ScheduledTasks []string `json:"scheduledTasks"`
}

// InstanceExecutions is a list of all executions assigned to an instance.
type InstanceExecutions struct {
	Executions []string `json:"executions"`
}

// InstanceStatus describes the availability of a instance.
type InstanceStatus string

const (
	// Starting instances are in the process of booting.
	Starting InstanceStatus = "starting"

	// Running instances are online.
	Running InstanceStatus = "running"

	// Stopping instances are shutting down.
	Stopping InstanceStatus = "stopping"

	// Stopped instances have been shut down, but not terminated.
	Stopped InstanceStatus = "stopped"

	// Terminated instances have been permanently stopped (deleted).
	Terminated InstanceStatus = "terminated"
)

// Deprecated. Use clusters/instances instead.
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
