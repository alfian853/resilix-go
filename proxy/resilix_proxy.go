package proxy

import (
	"resilix-go/context"
	"resilix-go/statehandler"
	"resilix-go/util"
)

type ResilixProxy struct {
	util.CheckedExecutor
	util.Executor
	statehandler.StateContainer

	stateHandler statehandler.StateHandler
}

func (proxy *ResilixProxy) Decorate(ctx *context.Context) *ResilixProxy {
	proxy.stateHandler = new(statehandler.CloseStateHandler).Decorate(ctx, proxy)

	return proxy
}

// StateContainer
func (proxy *ResilixProxy) SetStateHandler(stateHandler statehandler.StateHandler)  {
	proxy.stateHandler = stateHandler
}

func (proxy *ResilixProxy) GetStateHandler() statehandler.StateHandler {
	proxy.stateHandler.EvaluateState()
	return proxy.stateHandler
}

// CheckedExecutor
func (proxy *ResilixProxy) ExecuteChecked(fun func() error) (bool, error) {
	proxy.stateHandler.EvaluateState()

	return proxy.stateHandler.ExecuteChecked(fun)
}


func (proxy *ResilixProxy) ExecuteCheckedSupplier(fun func()(interface{}, error)) (bool, interface{}, error) {
	proxy.stateHandler.EvaluateState()

	return proxy.stateHandler.ExecuteCheckedSupplier(fun)
}


// Executor
func (proxy *ResilixProxy) Execute(fun func()) bool {

	isExecuted, err := proxy.ExecuteChecked(func() (err error) {
		fun()
		return nil
	})

	if err != nil {
		panic(err)
	}
	return isExecuted
}

func (proxy *ResilixProxy) ExecuteSupplier(fun func() interface{}) (bool, interface{})  {

	isExecuted, result, err := proxy.ExecuteCheckedSupplier(func() (interface{},error) {
		return fun(), nil
	})

	if err != nil {
		panic(err)
	}

	return isExecuted, result
}


