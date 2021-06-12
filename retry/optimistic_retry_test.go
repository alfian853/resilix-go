package retry

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
)

func TestOptimisticRetryRejected(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.3
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	var wg sync.WaitGroup

	retryExecutor := new(OptimisticRetryExecutor).Decorate(ctx)

	assert.Equal(t, RetryState_OnGoing, retryExecutor.GetRetryState())
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())

	minFailAck := int(ctx.Config.ErrorThreshold * float32(ctx.Config.NumberOfRetryInHalfOpenState))
	maxSuccessAck := int(ctx.Config.NumberOfRetryInHalfOpenState) - minFailAck

	for i := 0; i < maxSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed, result, err := retryExecutor.ExecuteCheckedSupplier(testutil.TrueCheckedSupplier())
			assert.True(t, executed)
			assert.True(t, result.(bool))
			assert.Nil(t, err)
		}, &wg)
	}

	for i := 0; i < minFailAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqError := testutil.RandErrorWithMessage()
			executed, err := retryExecutor.ExecuteChecked(testutil.PanicCheckedRunnable(uniqError))
			assert.True(t, executed)
			assert.Contains(t, err.Error(), uniqError.Error())
		}, &wg)
	}

	wg.Wait()

	assert.True(t, retryExecutor.getErrorRate() >= ctx.Config.ErrorThreshold)
	assert.False(t, retryExecutor.AcquirePermission())
	assert.Equal(t, RetryState_Rejected, retryExecutor.GetRetryState())
}

func TestOptimisticRetryAcceptedCase(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.8
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	var wg sync.WaitGroup

	retryExecutor := new(OptimisticRetryExecutor).Decorate(ctx)

	assert.Equal(t, RetryState_OnGoing, retryExecutor.GetRetryState())
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())

	minSuccessAck := int(math.Ceil(float64((1-ctx.Config.ErrorThreshold)*float32(ctx.Config.NumberOfRetryInHalfOpenState))) + 1)

	for i := 0; i < minSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed, err := retryExecutor.ExecuteChecked(testutil.CheckedRunnable())
			assert.True(t, executed)
			assert.Nil(t, err)
		}, &wg)
	}
	wg.Wait()
	for i := 0; i < int(ctx.Config.NumberOfRetryInHalfOpenState)-minSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			randErrorMessage := testutil.RandPanicMessage()
			executed, result, err := retryExecutor.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(randErrorMessage))
			assert.True(t, executed)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), randErrorMessage)
		}, &wg)
	}

	wg.Wait()

	assert.True(t, retryExecutor.getErrorRate() < ctx.Config.ErrorThreshold)
	assert.False(t, retryExecutor.AcquirePermission())
	assert.Equal(t, RetryState_Accepted, retryExecutor.GetRetryState())
}
