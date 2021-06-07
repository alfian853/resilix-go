package statehandler

import (
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/executor"
	"github.com/alfian853/resilix-go/slidingwindow"
)

type StateHandler interface {
	executor.CheckedExecutor
	EvaluateState()
}

type DefaultStateHandlerExt interface {
	acquirePermission() bool
	isSlidingWindowEnabled() bool
}

type DefaultStateHandler struct {
	StateHandler
	executor.DefaultExecutorExt

	defExecutor *executor.DefaultExecutor
	stateHandlerExt DefaultStateHandlerExt
	context *context.Context
	slidingWindow slidingwindow.SlidingWindow
}

func (defHandler *DefaultStateHandler) Decorate(
	ctx *context.Context, concreteHandler StateHandler, stateHandlerExt DefaultStateHandlerExt) *DefaultStateHandler {

	defHandler.defExecutor = new(executor.DefaultExecutor).Decorate(defHandler)
	defHandler.StateHandler = concreteHandler
	defHandler.context = ctx
	defHandler.slidingWindow = ctx.SWindow
	defHandler.stateHandlerExt = stateHandlerExt

	return defHandler
}

func (defHandler *DefaultStateHandler) ExecuteChecked(fun func() error) (executed bool, err error) {
	return defHandler.defExecutor.ExecuteChecked(fun)
}

func (defHandler *DefaultStateHandler) ExecuteCheckedSupplier(fun func()(interface{}, error)) (
	executed bool, result interface{}, err error) {
	return defHandler.defExecutor.ExecuteCheckedSupplier(fun)
}

func (defHandler *DefaultStateHandler) AcquirePermission() bool {
	return defHandler.stateHandlerExt.acquirePermission()
}

func (defHandler *DefaultStateHandler) OnAfterExecution(success bool) {
	defHandler.context.SWindow.AckAttempt(success)
	defHandler.StateHandler.EvaluateState()
}
