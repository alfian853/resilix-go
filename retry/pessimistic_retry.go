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

func (retry *PessimisticRetryExecutor) Decorate(ctx *context.Context) *PessimisticRetryExecutor {
	retry.ctx = ctx
	retry.config = ctx.Config
	retry.isAvailable = util.NewInt32(Available)
	retry.OptimisticRetryExecutor.DecorateWithSource(ctx, retry)
	return retry
}

func(retry *PessimisticRetryExecutor) acquireAndUpdateRetryPermission() bool {
	return atomic.SwapInt32(retry.isAvailable, NotAvailable) == Available &&
		retry.OptimisticRetryExecutor.acquireAndUpdateRetryPermission()
}

func (retry *PessimisticRetryExecutor) onAfterExecution() {
	atomic.SwapInt32(retry.isAvailable, Available)
}