package statehandler

import (
	conf "github.com/alfian853/resilix-go/config"
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

	defExecutor *executor.DefaultExecutor
	stateHandler StateHandler
	stateHandlerExt DefaultStateHandlerExt
	stateContainer StateContainer
	context *context.Context
	slidingWindow slidingwindow.SlidingWindow
	configuration *conf.Configuration
}

func (defHandler *DefaultStateHandler) Decorate(
	ctx *context.Context, concreteHandler StateHandler,
	stateHandlerExt DefaultStateHandlerExt, stateContainer StateContainer) *DefaultStateHandler {
	defHandler.defExecutor = new(executor.DefaultExecutor).Decorate(defHandler)
	defHandler.context = ctx
	defHandler.slidingWindow = ctx.SWindow
	defHandler.configuration = ctx.Config
	defHandler.stateContainer = stateContainer
	defHandler.stateHandlerExt = stateHandlerExt
	defHandler.stateHandler = concreteHandler
	defHandler.slidingWindow.SetActive(defHandler.stateHandlerExt.isSlidingWindowEnabled())

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
	defHandler.stateHandler.EvaluateState()
}
