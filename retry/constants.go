package retry

type RetryState int32

const (
	// RetryState
	RetryState_OnGoing  RetryState = 0
	RetryState_Rejected RetryState = 1
	RetryState_Accepted RetryState = 2
)

func NewRetryState(state RetryState) *RetryState {
	return &state
}
