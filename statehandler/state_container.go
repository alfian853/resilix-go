package statehandler

type StateContainer interface {
	setStateHandler(stateHandler StateHandler)
	getStateHandler() StateHandler
}