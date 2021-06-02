package context

import (
	"resilix-go/config"
	"resilix-go/slidingwindow"
)

type Context struct {
	Config  *config.Configuration
	SWindow slidingwindow.SlidingWindow
}

func NewContextDefault() *Context {
	ctx := Context{}
	ctx.Config = config.NewConfiguration()
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)

	return &ctx
}
