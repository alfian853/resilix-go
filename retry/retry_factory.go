package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
)

func CreateRetryExecutor(ctx *context.Context) RetryExecutor {
	switch ctx.Config.RetryStrategy {
	case consts.RETRY_OPTIMISTIC:
		return NewOptimisticRetryExecutor().Decorate(ctx)
	case consts.RETRY_PESSIMISTIC:
		return NewPessimisticRetryExecutor().Decorate(ctx)
	}

	return nil
}
