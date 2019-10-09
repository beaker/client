package api

import (
	"time"

	"github.com/shopspring/decimal"
)

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

	// InstanceCost describes the cost per node in units of USD-per-hour.
	InstanceCost *decimal.Decimal `json:"instanceCost,omitempty"`

	// Requested and actual configuration
	Status        ClusterStatus `json:"status"`
	InstanceSpec  InstanceSpec  `json:"instanceSpec"`
	InstanceShape *InstanceSpec `json:"instanceShape,omitempty"`
}

// ClusterPatch allows a client to update aspects of a Cluster.
type ClusterPatch struct {
	// Capacity changes the maximum number of instances a cluster can contain at one time.
	Capacity *int `json:"capacity,omitempty"`

	// Valid permanently sets validity for the cluster and should be accompanied
	// by an instance spec in the same request. If set to true, the cluster is
	// ready for use. Otherwise, it's considered failed.
	//
	// This value is internal; behavior is undefined if set by external clients.
	Valid *bool `json:"valid,omitempty"`

	// InstanceShape details the shape of instances created during cluster creation.
	//
	// This value is internal; behavior is undefined if set by external clients.
	InstanceShape *InstanceSpec `json:"instanceShape,omitempty"`

	// InstanceCost sets the estimated cost of each instance within the cluster in USD-per-hour.
	InstanceCost *decimal.Decimal `json:"instanceCost,omitempty"`
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

// InstanceSpec provides options to configure compute instances.
type InstanceSpec struct {
	CPUCount int    `json:"cpuCount,omitempty"`
	GPUCount int    `json:"gpuCount,omitempty"`
	GPUType  string `json:"gpuType,omitempty"`
	Memory   string `json:"memory,omitempty"`
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

type Instance struct {
	ID       string         `json:"id"`
	Hostname string         `json:"hostname"`
	Created  time.Time      `json:"created"`
	Status   InstanceStatus `json:"status"`
}

type ClusterInstances struct {
	// instances running on this cluster
	Data []Instance `json:"data"`
}

type CreateInstanceSpec struct {
	Hostname string `json:"hostname"`
}

type InstanceStatusSpec struct {
	Status InstanceStatus `json:"status"`
}

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

// ScheduledTasks is an ordered collection of task executions.
type ScheduledTasks struct {
	Data []ScheduledTask `json:"data"`
}

// ScheduledTask summarizes relations of tasks which have been scheduled.
type ScheduledTask struct {
	TaskID      string `json:"taskId"`
	ExecutionID string `json:"executionId"`

	// InstanceID is set when a task has been scheduled on an instance.
	InstanceID string `json:"instanceId,omitempty"`
}
