package retry

import "resilix-go/context"

func CreateRetryManager(ctx *context.Context) RetryManager {
	switch ctx.Configuration.RetryStrategy {
	case RETRY_OPTIMISTIC:
		return NewOptimisticRetryManager().Decorate(ctx)
	case RETRY_PESSIMISTIC:
		return NewPessimisticRetryManager().Decorate(ctx)
	}

	return nil
}
