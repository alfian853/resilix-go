package statehandler

type ResilixProxy struct {
	PanicExecutor
	Executor
	StateContainer
	stateHandler StateHandler
}

// StateContainer
func (proxy *ResilixProxy) setStateHandler(stateHandler StateHandler)  {
	proxy.stateHandler = stateHandler
}

func (proxy *ResilixProxy) getStateHandler() StateHandler  {
	return proxy.stateHandler
}


// PanicExecutor
func (proxy *ResilixProxy) PanicExecute(fun func() error) (bool, error) {
	proxy.stateHandler.evaluateState()

	return proxy.stateHandler.PanicExecute(fun)
}


func (proxy *ResilixProxy) PanicExecuteWithReturn(fun func()(interface{}, error)) (bool, interface{}, error) {
	proxy.stateHandler.evaluateState()

	return proxy.stateHandler.PanicExecuteWithReturn(fun)
}


// Executor
func (proxy *ResilixProxy) Execute(fun func()) bool {
	proxy.stateHandler.evaluateState()

	isExecuted, err := proxy.PanicExecute(func() (err error) {
		fun()
		return nil
	})

	if err != nil {
		panic(err)
	}
	return isExecuted
}

func (proxy *ResilixProxy) ExecuteWithReturn(fun func() interface{}) (bool, interface{})  {
	proxy.stateHandler.evaluateState()

	isExecuted, result, err := proxy.PanicExecuteWithReturn(func() (interface{},error) {
		return fun(), nil
	})

	if err != nil {
		panic(err)
	}

	return isExecuted, result
}


