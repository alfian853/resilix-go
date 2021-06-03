package retry

import (
	conf "resilix-go/config"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/util"
	"sync/atomic"
)

type OptimisticRetryManager struct {
	RetryManager
	slidingwindow.SwObserver

	numberOfRetry *int32
	numberOfFail  *int32
	ctx           *context.Context
	config        *conf.Configuration
}

func NewOptimisticRetryManager() *OptimisticRetryManager {return &OptimisticRetryManager{}}

func (retryManager *OptimisticRetryManager) Decorate(ctx *context.Context) *OptimisticRetryManager {
	retryManager.ctx = ctx
	retryManager.config = ctx.Config
	retryManager.numberOfRetry = util.NewInt32(0)
	retryManager.numberOfFail = util.NewInt32(0)
	ctx.SWindow.AddObserver(retryManager)
	return retryManager
}

func (retryManager *OptimisticRetryManager) AcquireAndUpdateRetryPermission() bool {

	numberOfRetry := atomic.AddInt32(retryManager.numberOfRetry,1) - 1

	if numberOfRetry >= retryManager.config.NumberOfRetryInHalfOpenState {
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

func (retryManager *OptimisticRetryManager) GetRetryState() consts.RetryState {
	if retryManager.isErrorLimitExceeded() {
		retryManager.ctx.SWindow.RemoveObserver(retryManager)
		return consts.RETRY_REJECTED
	}

	if atomic.LoadInt32(retryManager.numberOfRetry) >= retryManager.config.NumberOfRetryInHalfOpenState {
		retryManager.ctx.SWindow.RemoveObserver(retryManager)
		return consts.RETRY_ACCEPTED
	}

	return consts.RETRY_ON_GOING
}

func (retryManager *OptimisticRetryManager) NotifyOnAckAttempt(success bool) {
	if !success {
		atomic.AddInt32(retryManager.numberOfFail, 1)
	}
}

func (retryManager *OptimisticRetryManager) isErrorLimitExceeded() bool {
	return retryManager.GetErrorRate() >= retryManager.config.ErrorThreshold
}
