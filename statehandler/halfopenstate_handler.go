package statehandler

import (
	"resilix-go/context"
	"resilix-go/retry"
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
	stateHandler.timeEnd = getTimeStamp() + stateHandler.configuration.WaitDurationInOpenState

	stateHandler.retryManager = retry.CreateRetryManager(ctx)
	return stateHandler
}

func (stateHandler *HalfOpenStateHandler) acquirePermission() bool {
	return stateHandler.retryManager.AcquireAndUpdateRetryPermission()
}

func (stateHandler *HalfOpenStateHandler) evaluateState() {
	switch stateHandler.retryManager.GetRetryState() {
	case retry.ACCEPTED:
		newStateHandler := NewCloseStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newStateHandler)
		break

	case retry.REJECTED:
		newStateHandler := NewOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newStateHandler)
		break

	case retry.ON_GOING:
		// do nothing
		break
	}
}

