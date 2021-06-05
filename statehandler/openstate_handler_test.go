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

	stateHandler := new(OpenStateHandler).Decorate(ctx, &container)

	container.SetStateHandler(stateHandler)

	assert.False(t, stateHandler.acquirePermission())
	assert.False(t, container.GetStateHandler().acquirePermission())
	assert.Equal(t, stateHandler, container.GetStateHandler())

	time.Sleep(time.Duration(ctx.Config.WaitDurationInOpenState/2) * time.Millisecond)

	assert.False(t, stateHandler.acquirePermission())
	assert.False(t, container.GetStateHandler().acquirePermission())
	assert.Equal(t, stateHandler, container.GetStateHandler())

	executed, result, err := stateHandler.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier("won't happen"))

	assert.False(t, executed)
	assert.Nil(t, result)
	assert.Nil(t, err)

	executed, err = stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable("won't happen"))

	assert.False(t, executed)
	assert.Nil(t, err)

	time.Sleep(time.Duration(ctx.Config.WaitDurationInOpenState/2) * time.Millisecond)

	container.GetStateHandler().EvaluateState()
	assert.NotEqual(t, stateHandler, container.GetStateHandler())

	assert.IsType(t, &HalfOpenStateHandler{}, container.GetStateHandler())
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

	stateHandler := new(OpenStateHandler).Decorate(ctx, &container)

	container.SetStateHandler(stateHandler)

	for i:=0; i < ctx.Config.SlidingWindowMaxSize; i++ {
		ctx.SWindow.AckAttempt(false)
	}

	assert.Equal(t, stateHandler, container.GetStateHandler())
	assert.Equal(t, initialError, stateHandler.slidingWindow.GetErrorRate())
}
