package statehandler

import (
	"resilix-go/context"
	"time"
)

type OpenStateHandler struct {
	DefaultStateHandler
	timeEnd int64
}

func NewOpenStateHandler() *OpenStateHandler {
	return &OpenStateHandler{}
}

func (stateHandler *OpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *OpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateContainer)
	stateHandler.timeEnd = getTimeStamp() + stateHandler.configuration.WaitDurationInOpenState

	return stateHandler
}

func (stateHandler *HalfOpenStateHandler) isSlidingWindowActive() bool {
	return false
}

func (stateHandler *OpenStateHandler) acquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() <= stateHandler.configuration.ErrorThreshold
}

func (stateHandler *OpenStateHandler) evaluateState() {

	if stateHandler.timeEnd <= getTimeStamp() {
		newHandler := NewHalfOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newHandler)
	}
}

func getTimeStamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
