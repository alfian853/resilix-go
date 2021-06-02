package retry

import (
	"resilix-go/context"
	"sync/atomic"
)

type Availability *int32

const (
	Available    = 1
	NotAvailable = 0
)

type PessimisticRetryManager struct {
	OptimisticRetryManager

	isAvailable Availability
}


func NewPessimisticRetryManager() *PessimisticRetryManager {return &PessimisticRetryManager{}}

func (retryManager *PessimisticRetryManager) Decorate(ctx *context.Context) *PessimisticRetryManager {
	retryManager.ctx = ctx
	retryManager.config = ctx.Config

	ctx.SWindow.AddObserver(retryManager)
	return retryManager
}


func(retryManager *PessimisticRetryManager) AcquireAndUpdateRetryPermission() bool {
	return atomic.SwapInt32(retryManager.isAvailable, NotAvailable) == Available &&
		retryManager.OptimisticRetryManager.AcquireAndUpdateRetryPermission()
}


func (retryManager *PessimisticRetryManager) NotifyOnAckAttempt(success bool) {
	retryManager.OptimisticRetryManager.NotifyOnAckAttempt(success)

	atomic.SwapInt32(retryManager.isAvailable, Available)
}