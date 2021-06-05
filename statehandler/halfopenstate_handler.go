package statehandler

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/retry"
	"github.com/alfian853/resilix-go/util"
)

type HalfOpenStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt
	timeEnd       int64
	retryExecutor retry.RetryExecutor
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

func (stateHandler *HalfOpenStateHandler) acquirePermission() bool {
	return false
}

func (stateHandler *HalfOpenStateHandler) EvaluateState() {
	switch stateHandler.retryExecutor.GetRetryState() {
	case consts.RETRY_ACCEPTED:
		newStateHandler := new(CloseStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newStateHandler)
		break

	case consts.RETRY_REJECTED:
		newStateHandler := new(OpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newStateHandler)
		break

	case consts.RETRY_ON_GOING:
		// do nothing
		break
	}
}

