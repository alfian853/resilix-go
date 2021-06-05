package config

import (
	"github.com/alfian853/resilix-go/consts"
)

const (
	SECONDS_IN_MS = 1000
)


type Configuration struct {
	SlidingWindowStrategy        consts.SwStrategy
	RetryStrategy                consts.RetryStrategy
	SlidingWindowTimeRange       int64
	SlidingWindowMaxSize         int
	MinimumCallToEvaluate        int
	ErrorThreshold               float32
	WaitDurationInOpenState      int64
	NumberOfRetryInHalfOpenState int32
}

func NewConfiguration() *Configuration {
	config := Configuration{}

	config.SlidingWindowStrategy = consts.SwStrategy_CountBased
	config.RetryStrategy = consts.RETRY_PESSIMISTIC
	config.SlidingWindowTimeRange = 10 * SECONDS_IN_MS
	config.SlidingWindowMaxSize = 20
	config.MinimumCallToEvaluate = 5
	config.ErrorThreshold = 0.5
	config.WaitDurationInOpenState = 15 * SECONDS_IN_MS
	config.NumberOfRetryInHalfOpenState = 10

	return &config
}
