package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"sync/atomic"
)

type OptimisticRetryExecutor struct {
	DefaultRetryExecutor
	RetryExecutorExt
	retryState *consts.RetryState
}

func (retry *OptimisticRetryExecutor) Decorate(ctx *context.Context) *OptimisticRetryExecutor {
	retry.DefaultRetryExecutor.Decorate(ctx, retry)
	retry.retryState = consts.NewRetryState(consts.RETRY_ON_GOING)
	return retry
}

func (retry *OptimisticRetryExecutor) DecorateWithSource(
	ctx *context.Context, source RetryExecutorExt) *OptimisticRetryExecutor {
	retry.DefaultRetryExecutor.Decorate(ctx, source)
	retry.retryState = consts.NewRetryState(consts.RETRY_ON_GOING)
	return retry
}

func (retry *OptimisticRetryExecutor) GetRetryState() consts.RetryState {

	if consts.RetryState(atomic.LoadInt32((*int32)(retry.retryState))) != consts.RETRY_ON_GOING {
		return *retry.retryState
	}

	if retry.isErrorLimitExceeded() {
		atomic.SwapInt32((*int32)(retry.retryState), int32(consts.RETRY_REJECTED))
		return consts.RETRY_REJECTED
	}

	if atomic.LoadInt32(retry.numberOfAck) >= retry.config.NumberOfRetryInHalfOpenState {
		atomic.SwapInt32((*int32)(retry.retryState), int32(consts.RETRY_ACCEPTED))
		return consts.RETRY_ACCEPTED
	}

	return consts.RETRY_ON_GOING
}

/*
unsafe to be exposed to the public due to no execution guarantee after this method call
*/
func (retry *OptimisticRetryExecutor) acquireAndUpdateRetryPermission() bool {
	numberOfRetry := atomic.AddInt32(retry.numberOfRetry, 1) - 1

	if numberOfRetry >= retry.config.NumberOfRetryInHalfOpenState {
		return false
	} else if retry.isErrorLimitExceeded() {
		return false
	}

	return true
}

func (retry *OptimisticRetryExecutor) onAfterExecution()  {
	// do nothing
}

func (retry *OptimisticRetryExecutor) isErrorLimitExceeded() bool {
	return retry.getErrorRate() >= retry.config.ErrorThreshold
}


func (retry *OptimisticRetryExecutor) getErrorRate() float32 {
	if atomic.LoadInt32(retry.numberOfFail) == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(retry.numberOfFail)) / float32(atomic.LoadInt32(retry.numberOfAck))
}
