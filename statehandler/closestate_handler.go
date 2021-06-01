package statehandler


type CloseStateHandler struct {
	DefaultStateHandler
}

func NewCloseStateHandler() *CloseStateHandler {
	return &CloseStateHandler{}
}

func (stateHandler *CloseStateHandler) acquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() <= stateHandler.configuration.ErrorThreshold
}

func (stateHandler *CloseStateHandler) evaluateState() {

	if !stateHandler.acquirePermission() {
		newHandler := NewOpenStateHandler().Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.setStateHandler(newHandler)
	}

}