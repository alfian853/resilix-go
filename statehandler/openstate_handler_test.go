package statehandler

import (
	"github.com/stretchr/testify/assert"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/testutil"
	"testing"
	"time"
)

func TestOpenState_movingStateAfterWaitingDurationIsSatisfied(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 200
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewOpenStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	assert.False(t, stateHandler.AcquirePermission())
	assert.False(t, container.getStateHandler().AcquirePermission())
	assert.Equal(t, stateHandler, container.getStateHandler())

	time.Sleep(time.Duration(ctx.Config.WaitDurationInOpenState/2) * time.Millisecond)

	assert.False(t, stateHandler.AcquirePermission())
	assert.False(t, container.getStateHandler().AcquirePermission())
	assert.Equal(t, stateHandler, container.getStateHandler())

	executed, result, err := stateHandler.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier("won't happen"))

	assert.False(t, executed)
	assert.Nil(t, result)
	assert.Nil(t, err)

	executed, err = stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable("won't happen"))

	assert.False(t, executed)
	assert.Nil(t, err)

	time.Sleep(time.Duration(ctx.Config.WaitDurationInOpenState/2) * time.Millisecond)

	container.getStateHandler().EvaluateState()
	assert.NotEqual(t, stateHandler, container.getStateHandler())

	assert.IsType(t, &HalfOpenStateHandler{}, container.getStateHandler())
}

func TestOpenState_shouldNotAck(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)

	for i:=0; i < ctx.Config.SlidingWindowMaxSize; i++ {
		ctx.SWindow.AckAttempt(true)
	}

	for ctx.SWindow.GetErrorRate() < ctx.SWindow.GetErrorRate() {
		ctx.SWindow.AckAttempt(false)
	}

	initialError := ctx.SWindow.GetErrorRate()

 	container := testStateContainer{}

	stateHandler := NewOpenStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	for i:=0; i < ctx.Config.SlidingWindowMaxSize; i++ {
		ctx.SWindow.AckAttempt(false)
	}

	assert.Equal(t, stateHandler, container.getStateHandler())
	assert.Equal(t, initialError, stateHandler.slidingWindow.GetErrorRate())
}
