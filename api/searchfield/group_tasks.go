package searchfield

type GroupTask string

const (
	GroupTaskID         GroupTask = "taskId"
	GroupTaskStatus     GroupTask = "taskStatus"
	GroupExperimentID   GroupTask = "experimentId"
	GroupExperimentName GroupTask = "experimentName"
)

func (gt GroupTask) String() string { return string(gt) }
