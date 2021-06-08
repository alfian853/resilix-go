package statehandler

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/retry"
	"github.com/alfian853/resilix-go/util"
)

type HalfOpenStateHandler struct {
	DefaultStateHandler

	configuration  *conf.Configuration
	stateContainer StateContainer
	timeEnd        int64
	retryExecutor  retry.RetryExecutor
}

func (stateHandler *HalfOpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *HalfOpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler)
	stateHandler.stateContainer = stateContainer
	stateHandler.configuration = ctx.Config
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState
	stateHandler.retryExecutor = retry.CreateRetryExecutor(ctx)
	return stateHandler
}

func (stateHandler *HalfOpenStateHandler) ExecuteChecked(fun func() error) (executed bool, err error) {
	defer func() {
		stateHandler.EvaluateState()
	}()
	return stateHandler.retryExecutor.ExecuteChecked(fun)
}

func (stateHandler *HalfOpenStateHandler) ExecuteCheckedSupplier(fun func() (interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		stateHandler.EvaluateState()
	}()
	return stateHandler.retryExecutor.ExecuteCheckedSupplier(fun)
}

func (stateHandler *HalfOpenStateHandler) AcquirePermission() bool {
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
