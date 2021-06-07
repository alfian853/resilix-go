package executor

import (
	"github.com/alfian853/resilix-go/util"
)

type Executor interface {
	Execute(fun func()) bool
	ExecuteSupplier(fun func() interface{}) (bool, interface{})
}

type CheckedExecutor interface {
	ExecuteChecked(fun func() error) (bool, error)
	ExecuteCheckedSupplier(fun func() (interface{}, error)) (bool, interface{}, error)
}

type DefaultExecutorExt interface {
	AcquirePermission() bool
	OnAfterExecution(success bool)
}

type DefaultExecutor struct {
	CheckedExecutor
	executorExt DefaultExecutorExt
}

func (defExecutor *DefaultExecutor) Decorate(defaultExecutorExt DefaultExecutorExt) *DefaultExecutor {
	defExecutor.executorExt = defaultExecutorExt

	return defExecutor
}

func (defExecutor *DefaultExecutor) ExecuteChecked(fun func() error) (executed bool, err error) {
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
		if executed {
			defExecutor.executorExt.OnAfterExecution(err == nil)
		}
	}()

	if !defExecutor.executorExt.AcquirePermission() {
		return false, nil
	}
	executed = true
	err = fun()

	return true, err
}

func (defExecutor *DefaultExecutor) ExecuteCheckedSupplier(fun func() (interface{}, error)) (
	executed bool, result interface{}, err error) {
	defer func() {
		if message := recover(); message != nil {
			err = &util.UnhandledError{Message: message}
		}
		if executed {
			defExecutor.executorExt.OnAfterExecution(err == nil)
		}
	}()

	if !defExecutor.executorExt.AcquirePermission() {
		return false, nil, nil
	}

	executed = true
	result, err = fun()

	return true, result, err
}
