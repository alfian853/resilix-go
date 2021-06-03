package consts

type RetryStrategy string
type RetryState int
const (
	// RetryStrategy
	RETRY_PESSIMISTIC = "pessimistic"
	RETRY_OPTIMISTIC = "optimistic"


	// RetryState
	RETRY_ON_GOING = 0
	RETRY_REJECTED = 1
	RETRY_ACCEPTED = 2
)

