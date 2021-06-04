package statehandler

import (
	"github.com/stretchr/testify/assert"
	"math"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/testutil"
	"resilix-go/util"
	"sync"
	"testing"
)

func TestCloseStateHandler_minCallToEvaluate(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewCloseStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	for i:=0; i < ctx.Config.MinimumCallToEvaluate-1; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqPanic := testutil.RandPanicMessage()
			executed,result,err := stateHandler.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(uniqPanic))
			assert.True(t, executed)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), uniqPanic)
			assert.True(t, stateHandler.AcquirePermission())
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.getStateHandler())

	uniqPanic := testutil.RandPanicMessage()
	executed,err := stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable(&testutil.IntendedPanic{Message: uniqPanic}))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), uniqPanic)

	assert.NotEqual(t, stateHandler, container.getStateHandler())
	assert.IsType(t, &OpenStateHandler{}, container.getStateHandler())
}

func TestCloseStateHandler_stillCloseAfterNumberOfAck(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewCloseStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	for i:=0; i < ctx.Config.SlidingWindowMaxSize; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed,err := stateHandler.ExecuteChecked(testutil.CheckedRunnable())
			assert.True(t, executed)
			assert.Nil(t, err)
		}, &wg)
	}
	wg.Wait()
	assert.True(t, stateHandler.AcquirePermission())
	assert.Equal(t, stateHandler, container.getStateHandler())

	safeErrorAttempt := int(math.Ceil(float64(ctx.Config.SlidingWindowMaxSize)*float64(ctx.Config.ErrorThreshold)) - 1)
	for i:=0; i < safeErrorAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqPanic := testutil.RandErrorWithMessage()
			executed,result,err := stateHandler.ExecuteCheckedSupplier(testutil.ErrorCheckedSupplier(&uniqPanic))
			assert.True(t, executed)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), uniqPanic.Message)
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.getStateHandler())
}


func TestCloseStateHandler_moveToOpenState(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewCloseStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	for i:=0; i < ctx.Config.SlidingWindowMaxSize; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed,err := stateHandler.ExecuteChecked(testutil.CheckedRunnable())
			assert.True(t, executed)
			assert.Nil(t, err)
		}, &wg)
	}
	wg.Wait()
	assert.True(t, stateHandler.AcquirePermission())
	assert.Equal(t, stateHandler, container.getStateHandler())


	errorAttempt := int(math.Ceil(float64(ctx.Config.SlidingWindowMaxSize)*float64(1 - ctx.Config.ErrorThreshold)))
	for i:=0; i < errorAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqStringer := testutil.RandStringer()
			executed,result,err := stateHandler.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(&uniqStringer))
			assert.True(t, executed)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), uniqStringer.Message)
		}, &wg)

	}
	wg.Wait()
	assert.NotEqual(t, stateHandler, container.getStateHandler())
	assert.IsType(t, &OpenStateHandler{}, container.getStateHandler())
	assert.False(t, container.getStateHandler().AcquirePermission())
}

