package retry

import (
	"github.com/alfian853/resilix-go/consts"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func tryHardSuccess(t *testing.T, retryExecutor RetryExecutor, level int, isRecursive *bool) {
	if level > 20 {
		return
	}
	executed, err := retryExecutor.ExecuteChecked(testutil.CheckedRunnable())
	if !executed {
		*isRecursive = true
		testutil.RandSleep(1, 5)
		tryHardSuccess(t, retryExecutor, level+1, isRecursive)
		return
	}
	assert.True(t, executed)
	assert.Nil(t, err)
}

func tryHardFailed(t *testing.T, retryExecutor RetryExecutor, level int, isRecursive *bool) {
	if level > 20 {
		return
	}
	randErrorMessage := testutil.RandPanicMessage()
	executed, result, err := retryExecutor.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(randErrorMessage))
	if !executed {
		*isRecursive = true
		testutil.RandSleep(1, 5)
		tryHardFailed(t, retryExecutor, level+1, isRecursive)
		return
	}
	assert.True(t, executed)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), randErrorMessage)
}


func TestPessimisticRetryRejected(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.3
	ctx.Config.NumberOfRetryInHalfOpenState = 1000
	ctx.Config.RetryStrategy = consts.RetryStrategy_Pessimistic
	var wg sync.WaitGroup

	retryExecutor := new(PessimisticRetryExecutor).Decorate(ctx)

	assert.Equal(t, consts.RetryState_OnGoing, retryExecutor.GetRetryState())
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())

	minFailAck := 300
	maxSuccessAck := 700
	hasRecursiveCall := false
	for i := 0; i < maxSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardSuccess(t, retryExecutor, 1, &hasRecursiveCall)
		}, &wg)
	}
	wg.Wait()

	for i := 0; i < minFailAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardFailed(t, retryExecutor, 1, &hasRecursiveCall)
		}, &wg)
	}
	wg.Wait()
	assert.True(t, hasRecursiveCall)
	assert.True(t, retryExecutor.getErrorRate() >= ctx.Config.ErrorThreshold)
	assert.False(t, retryExecutor.AcquirePermission())
	assert.Equal(t, consts.RetryState_Rejected, retryExecutor.GetRetryState())
	assert.Equal(t, ctx.Config.NumberOfRetryInHalfOpenState, *retryExecutor.numberOfAck)
}

func TestPessimisticRetryAccepted(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.8
	ctx.Config.NumberOfRetryInHalfOpenState = 1000
	ctx.Config.RetryStrategy = consts.RetryStrategy_Pessimistic
	var wg sync.WaitGroup

	retryExecutor := new(PessimisticRetryExecutor).Decorate(ctx)

	assert.Equal(t, consts.RetryState_OnGoing, retryExecutor.GetRetryState())
	assert.Equal(t, float32(0), retryExecutor.getErrorRate())

	minSuccessAck := 801
	maxFailureAck := 199
	hasRecursiveCall := false

	for i := 0; i < minSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardSuccess(t, retryExecutor, 1, &hasRecursiveCall)
		}, &wg)
	}
	wg.Wait()
	for i := 0; i < maxFailureAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			tryHardFailed(t, retryExecutor, 1, &hasRecursiveCall)
		}, &wg)
	}

	wg.Wait()

	assert.True(t, hasRecursiveCall)
	assert.True(t, retryExecutor.getErrorRate() < ctx.Config.ErrorThreshold)
	assert.False(t, retryExecutor.AcquirePermission())
	assert.Equal(t, consts.RetryState_Accepted, retryExecutor.GetRetryState())
}
