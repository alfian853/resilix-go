package statehandler

import (
	conf "resilix-go/config"
	"resilix-go/context"
	"resilix-go/slidingwindow"
)

type PanicExecutor interface {
	PanicExecute(fun func() error) (bool, error)
	PanicExecuteWithReturn(fun func()(interface{}, error)) (bool, interface{}, error)
}

type StateHandler interface {
	PanicExecutor
	evaluateState()
	acquirePermission() bool
}

type DefaultStateHandler struct {
	StateHandler
	stateContainer StateContainer
	context *context.Context
	slidingWindow slidingwindow.SlidingWindow
	configuration *conf.Configuration
}

func (stateHandler *DefaultStateHandler) isSlidingWindowActive() bool {
	return true
}

func (stateHandler *DefaultStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *DefaultStateHandler {
	stateHandler.context = ctx
	stateHandler.slidingWindow = ctx.SWindow
	stateHandler.configuration = ctx.Config
	stateHandler.stateContainer = stateContainer

	ctx.SWindow.SetActive(stateHandler.isSlidingWindowActive())

	return stateHandler
}

func (stateHandler *DefaultStateHandler) PanicExecute(fun func() error) (bool, error) {

	if !stateHandler.acquirePermission() {
		return false, nil
	}
	err := fun()

	stateHandler.context.SWindow.AckAttempt(err == nil)
	stateHandler.evaluateState()

	return true, err
}

func (stateHandler *DefaultStateHandler) PanicExecuteWithReturn(fun func()(interface{}, error)) (bool, interface{}, error) {

	if !stateHandler.acquirePermission() {
		return false, nil, nil
	}
	result, err := fun()

	stateHandler.context.SWindow.AckAttempt(err == nil)
	stateHandler.evaluateState()

	return true, result, err
}
