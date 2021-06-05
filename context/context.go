package context

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/slidingwindow"
)

type Context struct {
	Config  *config.Configuration
	SWindow slidingwindow.SlidingWindow
}

func NewContextDefault() *Context {
	ctx := new(Context)
	ctx.Config = config.NewConfiguration()
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)

	return ctx
}

func NewContext(conf *config.Configuration) *Context {
	if conf == nil {
		return NewContextDefault()
	}

	ctx := new(Context)
	ctx.Config = conf
	ctx.SWindow = slidingwindow.CreateSlidingWindow(conf)

	return ctx
}
