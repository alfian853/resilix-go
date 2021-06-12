package retry

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/executor"
	"github.com/alfian853/resilix-go/util"
	"sync"
	"sync/atomic"
)

type OptimisticRetryExecutor struct {
	RetryExecutor
	executor.DefaultExecutorExt
	defExecutor   *executor.DefaultExecutor
	retryState    *RetryState
	writeNORLock  sync.Mutex
	numberOfRetry *int32
	numberOfAck   *int32
	numberOfFail  *int32
	ctx           *context.Context
	cfg           *config.Configuration
}

func (retry *OptimisticRetryExecutor) Decorate(ctx *context.Context) *OptimisticRetryExecutor {
	retry.defExecutor = new(executor.DefaultExecutor).Decorate(retry)
	retry.retryState = NewRetryState(RetryState_OnGoing)
	retry.ctx = ctx
	retry.cfg = ctx.Config
	retry.numberOfRetry = util.NewInt32(0)
	retry.numberOfFail = util.NewInt32(0)
	retry.numberOfAck = util.NewInt32(0)
	return retry
}

func (retry *OptimisticRetryExecutor) GetRetryState() RetryState {

	if RetryState(atomic.LoadInt32((*int32)(retry.retryState))) != RetryState_OnGoing {
		return *retry.retryState
	}

	if retry.isErrorLimitExceeded() {
		atomic.SwapInt32((*int32)(retry.retryState), int32(RetryState_Rejected))
		return RetryState_Rejected
	}

	if atomic.LoadInt32(retry.numberOfAck) >= retry.cfg.NumberOfRetryInHalfOpenState {
		atomic.SwapInt32((*int32)(retry.retryState), int32(RetryState_Accepted))
		return RetryState_Accepted
	}

	return RetryState_OnGoing
}

func (retry *OptimisticRetryExecutor) ExecuteChecked(fun func() error) (executed bool, err error) {
	return retry.defExecutor.ExecuteChecked(fun)
}

func (retry *OptimisticRetryExecutor) ExecuteCheckedSupplier(fun func() (interface{}, error)) (
	executed bool, result interface{}, err error) {
	return retry.defExecutor.ExecuteCheckedSupplier(fun)
}

/*
unsafe to be exposed to the public due to no execution guarantee after this method call
*/
func (retry *OptimisticRetryExecutor) AcquirePermission() bool {
	numberOfRetry := atomic.AddInt32(retry.numberOfRetry, 1) - 1

	if numberOfRetry >= retry.cfg.NumberOfRetryInHalfOpenState {
		return false
	} else if retry.isErrorLimitExceeded() {
		return false
	}

	return true
}

func (retry *OptimisticRetryExecutor) OnAfterExecution(success bool) {

	atomic.AddInt32(retry.numberOfAck, 1)
	if !success {
		atomic.AddInt32(retry.numberOfFail, 1)
	}
}

func (retry *OptimisticRetryExecutor) isErrorLimitExceeded() bool {
	return retry.getErrorRate() >= retry.cfg.ErrorThreshold
}

func (retry *OptimisticRetryExecutor) getErrorRate() float32 {
	if atomic.LoadInt32(retry.numberOfFail) == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(retry.numberOfFail)) / float32(atomic.LoadInt32(retry.numberOfAck))
}
