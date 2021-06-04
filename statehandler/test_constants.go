package statehandler


type testStateContainer struct {
	StateContainer

	stateHandler StateHandler
}

func (container *testStateContainer) setStateHandler(handler StateHandler)  {
	container.stateHandler = handler
}

func (container *testStateContainer) getStateHandler() StateHandler {
	return container.stateHandler
}

