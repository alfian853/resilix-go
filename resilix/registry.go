package resilix

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/proxy"
)

var executorMap = make(map[string]proxy.ResilixExecutor)

func Go(contextKey string) proxy.ResilixExecutor {
	val, ok := executorMap[contextKey]

	if ok {
		return val
	}

	return Register(contextKey, config.NewConfiguration())
}

func Register(contextKey string, cfg *config.Configuration) proxy.ResilixExecutor {
	ctx := context.NewContext(cfg)
	executor := proxy.NewResilixProxy(ctx)

	executorMap[contextKey] = executor

	return executor
}
