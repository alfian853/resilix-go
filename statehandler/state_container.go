package statehandler

type StateContainer interface {
	SetStateHandler(stateHandler StateHandler)
	GetStateHandler() StateHandler
}
