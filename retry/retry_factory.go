package retry

import (
	"resilix-go/consts"
	"resilix-go/context"
)

func CreateRetryManager(ctx *context.Context) RetryManager {
	switch ctx.Config.RetryStrategy {
	case consts.RETRY_OPTIMISTIC:
		return NewOptimisticRetryManager().Decorate(ctx)
	case consts.RETRY_PESSIMISTIC:
		return NewPessimisticRetryManager().Decorate(ctx)
	}

	return nil
}
