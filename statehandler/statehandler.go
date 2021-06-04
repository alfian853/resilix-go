package statehandler

import (
	"fmt"
	conf "resilix-go/config"
	"resilix-go/context"
	"resilix-go/slidingwindow"
)

type CheckedExecutor interface {
	ExecuteChecked(fun func() error) (bool, error)
	ExecuteCheckedSupplier(fun func()(interface{}, error)) (bool, interface{}, error)
}

type UnhandledError struct {
	error
	message interface{}
}

func (e *UnhandledError) Error() string {

	preFormat := "ResilixExecutor encountered an unhandled error"

	switch e.message.(type) {
	case error:
		return fmt.Sprintf(preFormat + ": %s\n", e.message.(error).Error())
	case string:
		return fmt.Sprintf(preFormat + ": %s\n", e.message.(string))
	}

	var args []interface{}
	canBeString := false
	if _,ok := e.message.(fmt.Stringer); ok {
		canBeString = true
		preFormat += ", String(): %s"
		args = append(args, e.message.(fmt.Stringer).String())
	}

	if _,ok := e.message.(fmt.GoStringer); ok {
		canBeString = true
		preFormat += ", GoString(): %s"
		args = append(args, e.message.(fmt.Stringer).String())
	}

	if canBeString {
		return fmt.Sprintf(preFormat + "\n", args)
	}

	return fmt.Sprintf(preFormat + ", %%#v: %#v\n", e.message)
}

type StateHandler interface {
	CheckedExecutor
	EvaluateState()
	AcquirePermission() bool
}

type DefaultStateHandlerExt interface {
	isSlidingWindowEnabled() bool
}

type DefaultStateHandler struct {
	CheckedExecutor
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
			err = &UnhandledError{message: message}
		}
	}()

	if !defHandler.stateHandler.AcquirePermission() {
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
			err = &UnhandledError{message: message}
		}
	}()

	if !defHandler.stateHandler.AcquirePermission() {
		//fmt.Println(defHandler.context.SWindow.GetErrorRate())
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