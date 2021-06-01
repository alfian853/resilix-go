package retry

import (
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"sync/atomic"
)

type OptimisticRetryManager struct {
	RetryManager
	slidingwindow.SwObserver

	numberOfRetry *int32
	numberOfFail  *int32
	ctx           *context.Context
	configuration *context.Configuration
}

func NewOptimisticRetryManager() *OptimisticRetryManager {return &OptimisticRetryManager{}}

func (retryManager *OptimisticRetryManager) Decorate(ctx *context.Context) *OptimisticRetryManager {
	retryManager.ctx = ctx
	retryManager.configuration = &ctx.Configuration

	ctx.SlidingWindow.AddObserver(retryManager)
	return retryManager
}

func (retryManager *OptimisticRetryManager) AcquireAndUpdateRetryPermission() bool {

	numberOfRetry := atomic.AddInt32(retryManager.numberOfRetry,1) - 1

	if numberOfRetry >= retryManager.configuration.NumberOfRetryInHalfOpenState {
		return false
	}

	return !retryManager.isErrorLimitExceeded()
}

func (retryManager *OptimisticRetryManager) GetErrorRate() float32 {
	if *retryManager.numberOfFail == 0 {
		return 0
	}

	return float32(*retryManager.numberOfFail) / float32(*retryManager.numberOfRetry)
}

func (retryManager *OptimisticRetryManager) GetRetryState() RetryState {
	if retryManager.isErrorLimitExceeded() {
		retryManager.ctx.SlidingWindow.RemoveObserver(retryManager)
		return REJECTED
	}

	if atomic.LoadInt32(retryManager.numberOfRetry) >= retryManager.configuration.NumberOfRetryInHalfOpenState {
		retryManager.ctx.SlidingWindow.RemoveObserver(retryManager)
		return ACCEPTED
	}

	return ON_GOING
}

func (retryManager *OptimisticRetryManager) NotifyOnAckAttempt(success bool) {
	if !success {
		atomic.AddInt32(retryManager.numberOfFail, 1)
	}
}

func (retryManager *OptimisticRetryManager) isErrorLimitExceeded() bool {
	return retryManager.GetErrorRate() >= retryManager.configuration.ErrorThreshold
}
