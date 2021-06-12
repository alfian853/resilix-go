package statehandler

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
)

type CloseStateHandler struct {
	DefaultStateHandler

	cfg            *config.Configuration
	stateContainer StateContainer
}

func (stateHandler *CloseStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *CloseStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler)
	stateHandler.stateContainer = stateContainer
	stateHandler.cfg = ctx.Config
	stateHandler.slidingWindow.Clear()
	return stateHandler
}

func (stateHandler *CloseStateHandler) AcquirePermission() bool {
	return stateHandler.slidingWindow.GetErrorRate() < stateHandler.cfg.ErrorThreshold
}

func (stateHandler *CloseStateHandler) EvaluateState() {

	if !stateHandler.AcquirePermission() {
		newHandler := new(OpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}

}
