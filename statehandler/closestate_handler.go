package statehandler

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
)

type CloseStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt

	configuration  *conf.Configuration
	stateContainer StateContainer
}

func (stateHandler *CloseStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *CloseStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler)
	stateHandler.stateContainer = stateContainer
	stateHandler.configuration = ctx.Config
	stateHandler.slidingWindow.Clear()
	return stateHandler
}

func (stateHandler *CloseStateHandler) acquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() < stateHandler.configuration.ErrorThreshold
}

func (stateHandler *CloseStateHandler) EvaluateState() {

	if !stateHandler.AcquirePermission() {
		newHandler := new(OpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}

}
