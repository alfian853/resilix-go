package statehandler

import (
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/retry"
	"resilix-go/util"
)

type HalfOpenStateHandler struct {
	DefaultStateHandler
	timeEnd int64
	retryManager retry.RetryManager
}

func NewHalfOpenStateHandler() *HalfOpenStateHandler {
	return &HalfOpenStateHandler{}
}

func (stateHandler *HalfOpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *HalfOpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateContainer)
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState

	stateHandler.retryManager = retry.CreateRetryManager(ctx)
	return stateHandler
}

func (stateHandler *HalfOpenStateHandler) acquirePermission() bool {
	return stateHandler.retryManager.AcquireAndUpdateRetryPermission()
}

func (stateHandler *HalfOpenStateHandler) evaluateState() {
	switch stateHandler.retryManager.GetRetryState() {
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

