package schema

type State int64

const (
	StateInvalid State = iota
	StateFresh
	StateInProgress
	StateCrashed
	StateCompleted
)
