package config

type SwStrategy string

const (
	// sliding window type
	SwStrategy_CountBased SwStrategy = "count-based"
	SwStrategy_TimeBased  SwStrategy = "time-based"
)

type RetryStrategy string

const (
	// RetryStrategy
	RetryStrategy_Pessimistic RetryStrategy = "pessimistic"
	RetryStrategy_Optimistic  RetryStrategy = "optimistic"
)
