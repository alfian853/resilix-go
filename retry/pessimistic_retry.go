package retry

import (
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/util"
	"sync/atomic"
)

type Availability *int32

const (
	Available    = 1
	NotAvailable = 0
)

type PessimisticRetryExecutor struct {
	OptimisticRetryExecutor
	isAvailable Availability
}


func NewPessimisticRetryExecutor() *PessimisticRetryExecutor {return &PessimisticRetryExecutor{}}

func (retryExecutor *PessimisticRetryExecutor) Decorate(ctx *context.Context) *PessimisticRetryExecutor {
	retryExecutor.ctx = ctx
	retryExecutor.config = ctx.Config
	retryExecutor.isAvailable = util.NewInt32(Available)
	retryExecutor.OptimisticRetryExecutor.Decorate(ctx)
	return retryExecutor
}


func(retryExecutor *PessimisticRetryExecutor) acquireAndUpdateRetryPermission() bool {
	return atomic.SwapInt32(retryExecutor.isAvailable, NotAvailable) == Available &&
		retryExecutor.OptimisticRetryExecutor.acquireAndUpdateRetryPermission()
}