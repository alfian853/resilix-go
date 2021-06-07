package consts

type RetryStrategy string
type RetryState int32
const (
	// RetryStrategy
	RETRY_PESSIMISTIC = "pessimistic"
	RETRY_OPTIMISTIC = "optimistic"
)

const(
	// RetryState
	RETRY_ON_GOING RetryState = 0
	RETRY_REJECTED RetryState = 1
	RETRY_ACCEPTED RetryState = 2

)

func NewRetryState(state RetryState) *RetryState {
	return &state
}
