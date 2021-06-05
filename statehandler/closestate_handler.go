package statehandler

import "resilix-go/context"

type CloseStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt
}

func (stateHandler *CloseStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *CloseStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler, stateContainer)
	stateHandler.slidingWindow.Clear()
	return stateHandler
}

func (stateHandler *CloseStateHandler) isSlidingWindowEnabled() bool {
	return true
}

func (stateHandler *CloseStateHandler) acquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() < stateHandler.configuration.ErrorThreshold
}

func (stateHandler *CloseStateHandler) EvaluateState() {

	if !stateHandler.acquirePermission() {
		newHandler := new(OpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}

}