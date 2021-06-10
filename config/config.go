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
	config.RetryStrategy = consts.Retry_Pessimistic
	config.SlidingWindowTimeRange = 10 * SECONDS_IN_MS
	config.SlidingWindowMaxSize = 20
	config.MinimumCallToEvaluate = 5
	config.ErrorThreshold = 0.5
	config.WaitDurationInOpenState = 15 * SECONDS_IN_MS
	config.NumberOfRetryInHalfOpenState = 10

	return &config
}

func Validate(config *Configuration) {

	if !isValidSlidingWindowStrategy(config.SlidingWindowStrategy) {
		panic("SlidingWindowStrategy should be valid, please see the documentations")
	}

	if !isValidRetryStrategy(config.RetryStrategy) {
		panic("RetryStrategy should be valid, please see the documentation")
	}

	if config.SlidingWindowTimeRange <= 0 {
		panic("SlidingWindowTimeRange should be greater than 0")
	}

	if config.SlidingWindowMaxSize <= 0 {
		panic("SlidingWindowMaxSize should be greater than 0")
	}

	if config.ErrorThreshold <= 0 || config.ErrorThreshold > 1 {
		panic("ErrorThreshold should be between 0.0 and 1.0")
	}

	if config.WaitDurationInOpenState < 0 {
		panic("WaitDurationInOpenState should be greater or equal to 0")
	}

	if config.NumberOfRetryInHalfOpenState <= 0 {
		panic("NumberOfRetryInHalfOpenState should be greater than 0")
	}

	if config.MinimumCallToEvaluate < 0 {
		panic("MinimumCallToEvaluate should be greater than or equal to 0")
	}

}

func isValidSlidingWindowStrategy(strategy consts.SwStrategy) bool {
	switch strategy {
	case consts.SwStrategy_CountBased:
		return true
	case consts.SwStrategy_TimeBased:
		return true
	}

	return false
}

func isValidRetryStrategy(strategy consts.RetryStrategy) bool {
	switch strategy {
	case consts.Retry_Optimistic:
		return true
	case consts.Retry_Pessimistic:
		return true
	}

	return false
}