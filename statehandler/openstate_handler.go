package statehandler

import (
	"resilix-go/context"
	"resilix-go/util"
)

type OpenStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt
	timeEnd int64
}

func (stateHandler *OpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *OpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler, stateContainer)
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState

	return stateHandler
}

func (stateHandler *OpenStateHandler) isSlidingWindowEnabled() bool {
	return false
}

func (stateHandler *OpenStateHandler) acquirePermission() bool {
	stateHandler.EvaluateState()

	if stateHandler.stateContainer.GetStateHandler() != stateHandler {
		return stateHandler.stateContainer.GetStateHandler().acquirePermission()
	}
	return false
}

func (stateHandler *OpenStateHandler) EvaluateState() {

	if stateHandler.timeEnd <= util.GetTimestamp() {
		newHandler := new(HalfOpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}
}
