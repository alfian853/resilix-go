package retry

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/util"
	"sync"
	"sync/atomic"
)

type OptimisticRetryExecutor struct {
	RetryExecutor
	writeNORLock sync.Mutex
	numberOfRetry *int32
	numberOfAck   *int32
	numberOfFail  *int32
	ctx           *context.Context
	config        *conf.Configuration
}

func NewOptimisticRetryExecutor() *OptimisticRetryExecutor {return &OptimisticRetryExecutor{}}

func (retryExecutor *OptimisticRetryExecutor) Decorate(ctx *context.Context) *OptimisticRetryExecutor {
	retryExecutor.ctx = ctx
	retryExecutor.config = ctx.Config
	retryExecutor.numberOfRetry = util.NewInt32(0)
	retryExecutor.numberOfFail = util.NewInt32(0)
	retryExecutor.numberOfAck = util.NewInt32(0)
	return retryExecutor
}


func (retryExecutor *OptimisticRetryExecutor) ExecuteChecked(fun func() error) (executed bool, err error) {
	defer func() {
		if executed {
			atomic.AddInt32(retryExecutor.numberOfAck, 1)
			if err != nil {
				atomic.AddInt32(retryExecutor.numberOfFail, 1)
			}
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !retryExecutor.acquireAndUpdateRetryPermission() {
		return false, nil
	}
	executed = true
	err = fun()

	return true, err
}

func (retryExecutor *OptimisticRetryExecutor) ExecuteCheckedSupplier(fun func()(interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		if executed {
			atomic.AddInt32(retryExecutor.numberOfAck, 1)
			if err != nil {
				atomic.AddInt32(retryExecutor.numberOfFail, 1)
			}
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !retryExecutor.acquireAndUpdateRetryPermission() {
		return false, nil, nil
	}

	executed = true
	result, err = fun()

	return true, result, err
}

func (retryExecutor *OptimisticRetryExecutor) GetRetryState() consts.RetryState {
	if retryExecutor.isErrorLimitExceeded() {
		return consts.RETRY_REJECTED
	}

	if atomic.LoadInt32(retryExecutor.numberOfRetry) >= retryExecutor.config.NumberOfRetryInHalfOpenState &&
		atomic.LoadInt32(retryExecutor.numberOfAck) == retryExecutor.config.NumberOfRetryInHalfOpenState {
		return consts.RETRY_ACCEPTED
	}

	return consts.RETRY_ON_GOING
}

/*
unsafe to be exposed to the public due to no execution guarantee after this method call
*/
func (retryExecutor *OptimisticRetryExecutor) acquireAndUpdateRetryPermission() bool {

	numberOfRetry := atomic.AddInt32(retryExecutor.numberOfRetry, 1) - 1

	if numberOfRetry >= retryExecutor.config.NumberOfRetryInHalfOpenState {
		return false
	} else if retryExecutor.isErrorLimitExceeded() {
		return false
	}

	return true
}

func (retryExecutor *OptimisticRetryExecutor) isErrorLimitExceeded() bool {
	return retryExecutor.getErrorRate() >= retryExecutor.config.ErrorThreshold
}


func (retryExecutor *OptimisticRetryExecutor) getErrorRate() float32 {
	if atomic.LoadInt32(retryExecutor.numberOfFail) == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(retryExecutor.numberOfFail)) / float32(atomic.LoadInt32(retryExecutor.numberOfAck))
}
