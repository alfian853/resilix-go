package statehandler

type ResilixProxy struct {
	CheckedExecutor
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
	proxy.stateHandler.EvaluateState()

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
	proxy.stateHandler.EvaluateState()

	isExecuted, result, err := proxy.ExecuteCheckedSupplier(func() (interface{},error) {
		return fun(), nil
	})

	if err != nil {
		panic(err)
	}

	return isExecuted, result
}


