package retry

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetryFactory(t *testing.T) {
	ctx := context.NewContextDefault()

	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	assert.IsType(t, &OptimisticRetryExecutor{}, CreateRetryExecutor(ctx))

	ctx.Config.RetryStrategy = config.RetryStrategy_Pessimistic
	assert.IsType(t, &PessimisticRetryExecutor{}, CreateRetryExecutor(ctx))

	ctx.Config.RetryStrategy = "nothing"
	assert.PanicsWithValue(t, "Unknown RetryStrategy: nothing", func() {
		CreateRetryExecutor(ctx)
	})

}
