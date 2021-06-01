package context

import (
	"resilix-go/retry"
	"resilix-go/slidingwindow"
)

type Context struct {
	Configuration Configuration
	SlidingWindow slidingwindow.SlidingWindow
}

type Configuration struct {
	RetryStrategy                retry.Strategy
	SlidingWindowStrategy        slidingwindow.Strategy
	SlidingWindowTimeRange       uint
	SlidingWindowMaxSize         int
	MinimumCallToEvaluate        int
	ErrorThreshold               float32
	WaitDurationInOpenState      int64
	NumberOfRetryInHalfOpenState int32
}