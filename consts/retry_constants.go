package consts

type RetryStrategy string
type RetryState int
const (
	// RetryStrategy
	RETRY_PESSIMISTIC = "pessimistic"
	RETRY_OPTIMISTIC = "optimistic"


	// RetryState
	ON_GOING = 0
	REJECTED = 1
	ACCEPTED = 2
)

