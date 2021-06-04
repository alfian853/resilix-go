package retry

import (
	conf "resilix-go/config"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/util"
	"sync"
	"sync/atomic"
)

type OptimisticRetryManager struct {
	RetryManager
	slidingwindow.SwObserver
	writeNORLock sync.Mutex
	numberOfRetry *int32
	numberOfAck   *int32
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
	retryManager.numberOfAck = util.NewInt32(0)
	ctx.SWindow.AddObserver(retryManager)
	return retryManager
}

func (retryManager *OptimisticRetryManager) AcquireAndUpdateRetryPermission() bool {
	defer retryManager.writeNORLock.Unlock()
	retryManager.writeNORLock.Lock()

	allowed := !retryManager.isErrorLimitExceeded() &&
		atomic.LoadInt32(retryManager.numberOfRetry) < retryManager.config.NumberOfRetryInHalfOpenState

	if allowed {
		atomic.AddInt32(retryManager.numberOfRetry, 1)
	}

	return allowed
}

func (retryManager *OptimisticRetryManager) GetErrorRate() float32 {
	if atomic.LoadInt32(retryManager.numberOfFail) == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(retryManager.numberOfFail)) / float32(atomic.LoadInt32(retryManager.numberOfRetry))
}

func (retryManager *OptimisticRetryManager) GetRetryState() consts.RetryState {
	if retryManager.isErrorLimitExceeded() {
		retryManager.ctx.SWindow.RemoveObserver(retryManager)
		return consts.RETRY_REJECTED
	}

	if atomic.LoadInt32(retryManager.numberOfRetry) >= retryManager.config.NumberOfRetryInHalfOpenState &&
		atomic.LoadInt32(retryManager.numberOfAck) == atomic.LoadInt32(retryManager.numberOfRetry) {
		retryManager.ctx.SWindow.RemoveObserver(retryManager)
		return consts.RETRY_ACCEPTED
	}

	return consts.RETRY_ON_GOING
}

func (retryManager *OptimisticRetryManager) NotifyOnAckAttempt(success bool) {
	atomic.AddInt32(retryManager.numberOfAck, 1)
	if !success {
		atomic.AddInt32(retryManager.numberOfFail, 1)
	}
}

func (retryManager *OptimisticRetryManager) isErrorLimitExceeded() bool {
	return retryManager.GetErrorRate() >= retryManager.config.ErrorThreshold
}
