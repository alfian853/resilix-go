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

func NewOpenStateHandler() *OpenStateHandler {
	return &OpenStateHandler{}
}

func (stateHandler *OpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *OpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler, stateContainer)
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState

	return stateHandler
}

func (stateHandler *OpenStateHandler) isSlidingWindowEnabled() bool {
	return false
}

func (stateHandler *OpenStateHandler) AcquirePermission() bool {
	stateHandler.EvaluateState()

	if stateHandler.stateContainer.getStateHandler() != stateHandler {
		return stateHandler.stateContainer.getStateHandler().AcquirePermission()
	}
	return false
}

func (stateHandler *OpenStateHandler) EvaluateState() {

	if stateHandler.timeEnd <= util.GetTimestamp() {
		newHandler := NewHalfOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newHandler)
	}
}
