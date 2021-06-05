package statehandler


type testStateContainer struct {
	StateContainer

	stateHandler StateHandler
}

func (container *testStateContainer) SetStateHandler(handler StateHandler)  {
	container.stateHandler = handler
}

func (container *testStateContainer) GetStateHandler() StateHandler {
	return container.stateHandler
}

