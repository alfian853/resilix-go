package statehandler

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/util"
)

type OpenStateHandler struct {
	DefaultStateHandler
	DefaultStateHandlerExt

	configuration  *conf.Configuration
	stateContainer StateContainer
	timeEnd        int64
}

func (stateHandler *OpenStateHandler) Decorate(ctx *context.Context, stateContainer StateContainer) *OpenStateHandler {
	stateHandler.DefaultStateHandler.Decorate(ctx, stateHandler, stateHandler)
	stateHandler.stateContainer = stateContainer
	stateHandler.configuration = ctx.Config
	stateHandler.timeEnd = util.GetTimestamp() + stateHandler.configuration.WaitDurationInOpenState

	return stateHandler
}

func (stateHandler *OpenStateHandler) acquirePermission() bool {
	stateHandler.EvaluateState()
	return false
}

func (stateHandler *OpenStateHandler) EvaluateState() {

	if stateHandler.timeEnd <= util.GetTimestamp() {
		newHandler := new(HalfOpenStateHandler).Decorate(stateHandler.context, stateHandler.stateContainer)
		stateHandler.stateContainer.SetStateHandler(newHandler)
	}
}
