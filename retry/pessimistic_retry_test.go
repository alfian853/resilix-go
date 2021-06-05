package retry

import (
	"github.com/stretchr/testify/assert"
	"math"
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"sync"
	"testing"
)

func tryHardSuccess(t *testing.T, retryExecutor RetryExecutor) {
	executed, err := retryExecutor.ExecuteChecked(testutil.CheckedRunnable())
	if !executed {
		tryHardSuccess(t, retryExecutor)
	}
	assert.True(t, executed)
	assert.Nil(t, err)
}

func tryHardFailed(t *testing.T, retryExecutor RetryExecutor) {
	randErrorMessage := testutil.RandPanicMessage()
	executed, result, err := retryExecutor.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(randErrorMessage))
	if !executed {
		tryHardFailed(t, retryExecutor)
	}
	assert.True(t, executed)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), randErrorMessage)
}

func TestPessimisticRetryRejected(t *testing.T){
	t.Deadline()
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.3
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_PESSIMISTIC
	var wg sync.WaitGroup

	retryExecutor := NewPessimisticRetryExecutor().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryExecutor.GetRetryState()))
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())

	minFailAck := int(ctx.Config.ErrorThreshold * float32(ctx.Config.NumberOfRetryInHalfOpenState))
	maxSuccessAck := int(ctx.Config.NumberOfRetryInHalfOpenState) - minFailAck

	for i:=0 ; i < maxSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardSuccess(t, retryExecutor)
		}, &wg)
	}
	wg.Wait()

	for i:=0 ; i < minFailAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardFailed(t, retryExecutor)
		}, &wg)
	}
	wg.Wait()
	assert.True(t, retryExecutor.getErrorRate() >= ctx.Config.ErrorThreshold)
	assert.False(t, retryExecutor.acquireAndUpdateRetryPermission())
	assert.Equal(t, consts.RETRY_REJECTED, int(retryExecutor.GetRetryState()))
	assert.Equal(t, ctx.Config.NumberOfRetryInHalfOpenState, *retryExecutor.numberOfAck)
}

func TestPessimisticRetryAcceptedCase(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.8
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_PESSIMISTIC
	var wg sync.WaitGroup

	retryExecutor := NewPessimisticRetryExecutor().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryExecutor.GetRetryState()))
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())


	minSuccessAck := int(math.Ceil(float64((1 - ctx.Config.ErrorThreshold) * float32(ctx.Config.NumberOfRetryInHalfOpenState))) + 1)

	for i:=0 ; i < minSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardSuccess(t, retryExecutor)
		}, &wg)
	}
	wg.Wait()
	for i:=0 ; i < int(ctx.Config.NumberOfRetryInHalfOpenState) - minSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardFailed(t, retryExecutor)
		}, &wg)
	}

	wg.Wait()

	assert.True(t,  retryExecutor.getErrorRate() < ctx.Config.ErrorThreshold)
	assert.False(t,  retryExecutor.acquireAndUpdateRetryPermission())
	assert.Equal(t, consts.RETRY_ACCEPTED, int(retryExecutor.GetRetryState()))
}