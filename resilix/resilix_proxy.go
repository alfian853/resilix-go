package resilix

import (
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/executor"
	"github.com/alfian853/resilix-go/statehandler"
)

type ResilixExecutor interface {
	executor.CheckedExecutor
	Execute(fun func()) (bool, error)
	ExecuteSupplier(fun func() interface{}) (bool, interface{}, error)
}

type ResilixProxy struct {
	ResilixExecutor
	statehandler.StateContainer
	stateHandler statehandler.StateHandler
}

func NewResilixProxy(ctx *context.Context) *ResilixProxy {
	proxy := new(ResilixProxy)
	proxy.stateHandler = new(statehandler.CloseStateHandler).Decorate(ctx, proxy)

	return proxy
}

func (proxy *ResilixProxy) SetStateHandler(stateHandler statehandler.StateHandler) {
	proxy.stateHandler = stateHandler
}

func (proxy *ResilixProxy) GetStateHandler() statehandler.StateHandler {
	proxy.stateHandler.EvaluateState()
	return proxy.stateHandler
}

// ExecuteChecked returns isExecuted, error (any returned or runtime error)
func (proxy *ResilixProxy) ExecuteChecked(fun func() error) (bool, error) {
	proxy.stateHandler.EvaluateState()

	return proxy.stateHandler.ExecuteChecked(fun)
}

// ExecuteCheckedSupplier returns: isExecuted, result, error (any returned or runtime error)
func (proxy *ResilixProxy) ExecuteCheckedSupplier(fun func() (interface{}, error)) (bool, interface{}, error) {
	proxy.stateHandler.EvaluateState()

	return proxy.stateHandler.ExecuteCheckedSupplier(fun)
}

// Execute returns isExecuted, error (any runtime error)
func (proxy *ResilixProxy) Execute(fun func()) (bool, error) {

	isExecuted, err := proxy.ExecuteChecked(func() (err error) {
		fun()
		return nil
	})

	return isExecuted, err
}

// ExecuteSupplier returns: isExecuted, result, error (any runtime error)
func (proxy *ResilixProxy) ExecuteSupplier(fun func() interface{}) (bool, interface{}, error) {

	isExecuted, result, err := proxy.ExecuteCheckedSupplier(func() (interface{}, error) {
		return fun(), nil
	})

	return isExecuted, result, err
}
