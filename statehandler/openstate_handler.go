package statehandler

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/util"
)

type OpenStateHandler struct {
	DefaultStateHandler

	cfg            *config.Configuration
	stateContainer StateContainer
	timeEnd        int64
}

func (stateHandler *OpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *OpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler)
	stateHandler.stateContainer = stateContainer
	stateHandler.cfg = ctx.Config
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.cfg.WaitDurationInOpenState

	return stateHandler
}

func (stateHandler *OpenStateHandler) AcquirePermission() bool {
	stateHandler.EvaluateState()
	return false
}

func (stateHandler *OpenStateHandler) EvaluateState() {

	if stateHandler.timeEnd <= util.GetTimestamp() {
		newHandler := new(HalfOpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}
}
