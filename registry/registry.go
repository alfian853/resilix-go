package registry

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/proxy"
)

var	executorMap = make(map[string]proxy.ResilixExecutor)

func GetResilixExecutor(contextKey string) proxy.ResilixExecutor {
	val, ok := executorMap[contextKey]

	if ok {
		return val
	}

	return Register(contextKey, conf.NewConfiguration())
}

func Register(contextKey string, config *conf.Configuration) proxy.ResilixExecutor {
	ctx := context.NewContext(config)
	executor := proxy.NewResilixProxy(ctx)

	executorMap[contextKey] = executor

	return executor
}