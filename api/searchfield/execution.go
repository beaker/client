package searchfield

type Execution string

const (
	ExecutionID Execution = "id"
)

func (e Execution) String() string { return string(e) }
