package retry

import (
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/executor"
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
	retry.OptimisticRetryExecutor.Decorate(ctx)
	retry.defExecutor = new(executor.DefaultExecutor).Decorate(retry)
	retry.isAvailable = util.NewInt32(Available)
	return retry
}

func (retry *PessimisticRetryExecutor) AcquirePermission() bool {
	return atomic.SwapInt32(retry.isAvailable, NotAvailable) == Available &&
		retry.OptimisticRetryExecutor.AcquirePermission()
}

func (retry *PessimisticRetryExecutor) OnAfterExecution(success bool) {

	atomic.AddInt32(retry.numberOfAck, 1)
	if !success {
		atomic.AddInt32(retry.numberOfFail, 1)
	}
	atomic.SwapInt32(retry.isAvailable, Available)
}
