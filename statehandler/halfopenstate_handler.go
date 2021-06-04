package statehandler

import (
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/retry"
	"resilix-go/util"
)

type HalfOpenStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt
	timeEnd       int64
	retryExecutor retry.RetryExecutor
}

func NewHalfOpenStateHandler() *HalfOpenStateHandler {
	return &HalfOpenStateHandler{}
}

func (stateHandler *HalfOpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *HalfOpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler, stateContainer)
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState

	stateHandler.retryExecutor = retry.CreateRetryExecutor(ctx)
	return stateHandler
}

func(stateHandler *HalfOpenStateHandler) ExecuteChecked(fun func() error) (executed bool, err error) {
	defer func() {
		stateHandler.EvaluateState()
	}()
	return stateHandler.retryExecutor.ExecuteChecked(fun)
}

func (stateHandler *HalfOpenStateHandler) ExecuteCheckedSupplier(fun func()(interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		stateHandler.EvaluateState()
	}()
	return stateHandler.retryExecutor.ExecuteCheckedSupplier(fun)
}

func (stateHandler *HalfOpenStateHandler) isSlidingWindowEnabled() bool {
	return true
}

func (stateHandler *HalfOpenStateHandler) AcquirePermission() bool {
	return false
}

func (stateHandler *HalfOpenStateHandler) EvaluateState() {
	switch stateHandler.retryExecutor.GetRetryState() {
	case consts.RETRY_ACCEPTED:
		newStateHandler := NewCloseStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newStateHandler)
		break

	case consts.RETRY_REJECTED:
		newStateHandler := NewOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newStateHandler)
		break

	case consts.RETRY_ON_GOING:
		// do nothing
		break
	}
}

