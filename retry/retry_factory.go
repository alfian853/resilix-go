package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
)

func CreateRetryExecutor(ctx *context.Context) RetryExecutor {
	switch ctx.Config.RetryStrategy {
	case consts.Retry_Optimistic:
		return new(OptimisticRetryExecutor).Decorate(ctx)
	case consts.Retry_Pessimistic:
		return new(PessimisticRetryExecutor).Decorate(ctx)
	}

	panic("Unknown RetryStrategy: " + ctx.Config.RetryStrategy)
}
