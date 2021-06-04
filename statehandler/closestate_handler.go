package statehandler

import "resilix-go/context"

type CloseStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt
}

func NewCloseStateHandler() *CloseStateHandler {
	return &CloseStateHandler{}
}


func (stateHandler *CloseStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *CloseStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler, stateContainer)
	stateHandler.slidingWindow.Clear()
	return stateHandler
}

func (stateHandler *CloseStateHandler) isSlidingWindowEnabled() bool {
	return true
}

func (stateHandler *CloseStateHandler) AcquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() <= stateHandler.configuration.ErrorThreshold
}

func (stateHandler *CloseStateHandler) EvaluateState() {

	if !stateHandler.AcquirePermission() {
		newHandler := NewOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newHandler)
	}

}