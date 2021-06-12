package retry

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
)

func CreateRetryExecutor(ctx *context.Context) RetryExecutor {
	switch ctx.Config.RetryStrategy {
	case config.RetryStrategy_Optimistic:
		return new(OptimisticRetryExecutor).Decorate(ctx)
	case config.RetryStrategy_Pessimistic:
		return new(PessimisticRetryExecutor).Decorate(ctx)
	}

	panic("Unknown RetryStrategy: " + string(ctx.Config.RetryStrategy))
}
