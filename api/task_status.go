package api

// TaskStatus represents the stage of execution for a task. The associated enumeration is ordered,
// where higher value statuses are closer to complete. It's possible for a task to transition from a
// higher status to a lower one if rescheduled. For example, a "running" experiment can move to
// "initializing" if the node it's running on is terminated.
type TaskStatus string

const (
	// TaskStatusSubmitted means a task is accepted by Beaker.
	// The task will automatically start when eligible.
	TaskStatusSubmitted TaskStatus = "submitted"

	// TaskStatusProvisioning means a task has been submitted for execution and
	// Beaker is waiting for compute resources to become available.
	TaskStatusProvisioning TaskStatus = "provisioning"

	// TaskStatusInitializing means a task is attempting to start, but is not yet running.
	TaskStatusInitializing TaskStatus = "initializing"

	// TaskStatusRunning means a task is executing.
	TaskStatusRunning TaskStatus = "running"

	// TaskStatusTerminating means a task has finished executing,
	// and Beaker is processing the end results (e.g. logs, result data).
	TaskStatusTerminating TaskStatus = "terminating"

	// TaskStatusPreempted means a task was interrupted by an external cause, such as its node shutting down.
	TaskStatusPreempted TaskStatus = "preempted"

	// TaskStatusSucceeded means a task has completed successfully.
	TaskStatusSucceeded TaskStatus = "succeeded"

	// TaskStatusStopped means a task was interrupted.
	TaskStatusStopped TaskStatus = "stopped"

	// TaskStatusFailed means a task could not be completed.
	TaskStatusFailed TaskStatus = "failed"
)

// IsEndState is true if the TaskStatus is preempted, canceled, failed, or successful
func (ts TaskStatus) IsEndState() bool {
	switch ts {
	case TaskStatusPreempted, TaskStatusSucceeded, TaskStatusStopped, TaskStatusFailed:
		return true
	default:
		return false
	}
}
