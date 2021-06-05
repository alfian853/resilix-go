package retry

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/util"
	"sync"
	"sync/atomic"
)

type RetryExecutor interface {
	util.CheckedExecutor
	GetRetryState() consts.RetryState
}

type RetryExecutorExt interface {
	acquireAndUpdateRetryPermission() bool
	onAfterExecution()
}

type DefaultRetryExecutor struct {
	RetryExecutor
	retryExecutorExt RetryExecutorExt
	writeNORLock sync.Mutex
	numberOfRetry *int32
	numberOfAck   *int32
	numberOfFail  *int32
	ctx           *context.Context
	config        *conf.Configuration

}

func (retry *DefaultRetryExecutor) Decorate(ctx *context.Context, ext RetryExecutorExt) *DefaultRetryExecutor {
	retry.ctx = ctx
	retry.config = ctx.Config
	retry.numberOfRetry = util.NewInt32(0)
	retry.numberOfFail = util.NewInt32(0)
	retry.numberOfAck = util.NewInt32(0)
	retry.retryExecutorExt = ext
	return retry
}


func (retry *DefaultRetryExecutor) ExecuteChecked(fun func() error) (executed bool, err error) {
	defer func() {
		if executed {
			atomic.AddInt32(retry.numberOfAck, 1)
			if err != nil {
				atomic.AddInt32(retry.numberOfFail, 1)
			}
			retry.retryExecutorExt.onAfterExecution()
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !retry.retryExecutorExt.acquireAndUpdateRetryPermission() {
		return false, nil
	}
	executed = true
	err = fun()

	return true, err
}


func (retry *DefaultRetryExecutor) ExecuteCheckedSupplier(fun func()(interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		if executed {
			atomic.AddInt32(retry.numberOfAck, 1)
			if err != nil {
				atomic.AddInt32(retry.numberOfFail, 1)
			}
			retry.retryExecutorExt.onAfterExecution()
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !retry.retryExecutorExt.acquireAndUpdateRetryPermission() {
		return false, nil, nil
	}

	executed = true
	result, err = fun()

	return true, result, err
}
