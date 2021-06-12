package resilix

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
)

var executorMap = make(map[string]ResilixExecutor)

func Go(contextKey string) ResilixExecutor {
	val, ok := executorMap[contextKey]

	if ok {
		return val
	}

	return Register(contextKey, config.NewConfiguration())
}

func Register(contextKey string, cfg *config.Configuration) ResilixExecutor {
	ctx := context.NewContext(cfg)
	executor := NewResilixProxy(ctx)

	executorMap[contextKey] = executor

	return executor
}
