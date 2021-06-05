package statehandler

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/slidingwindow"
	"github.com/alfian853/resilix-go/util"
)

type StateHandler interface {
	util.CheckedExecutor
	EvaluateState()
	acquirePermission() bool
}

type DefaultStateHandlerExt interface {
	isSlidingWindowEnabled() bool
}

type DefaultStateHandler struct {
	StateHandler

	util.CheckedExecutor
	stateHandler StateHandler
	stateContainer StateContainer
	stateHandlerExt DefaultStateHandlerExt
	context *context.Context
	slidingWindow slidingwindow.SlidingWindow
	configuration *conf.Configuration
}

func (defHandler *DefaultStateHandler) Decorate(
	ctx *context.Context, concreteHandler StateHandler,
	stateHandlerExt DefaultStateHandlerExt, stateContainer StateContainer) *DefaultStateHandler {

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
	defer func() {
		if executed {
			defHandler.handleAfterExecution(err == nil)
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !defHandler.stateHandler.acquirePermission() {
		return false, nil
	}
	executed = true
	err = fun()

	return true, err
}

func (defHandler *DefaultStateHandler) ExecuteCheckedSupplier(fun func()(interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		if executed {
			defHandler.handleAfterExecution(err == nil)
		}
	}()
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
	}()

	if !defHandler.stateHandler.acquirePermission() {
		return false, nil, nil
	}

	executed = true
	result, err = fun()

	return true, result, err
}

func (defHandler *DefaultStateHandler)handleAfterExecution(success bool){
	defHandler.context.SWindow.AckAttempt(success)
	defHandler.stateHandler.EvaluateState()
}