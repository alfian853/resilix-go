package statehandler

import (
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/util"
	"testing"
)

func TestCloseStateHandler_minCallToEvaluate(t *testing.T) {
	//init
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
		uniqPanic := util.RandPanicMessage()
		executed,result,err := stateHandler.ExecuteCheckedSupplier(util.PanicCheckedSupplier(uniqPanic))
		assert.True(t, executed)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), uniqPanic)
		assert.True(t, stateHandler.AcquirePermission())
	}

	assert.Equal(t, stateHandler, container.getStateHandler())
	assert.True(
		t,
		reflect.TypeOf(container.getStateHandler()).AssignableTo(reflect.TypeOf(&CloseStateHandler{})),
	)

	uniqPanic := util.RandPanicMessage()
	executed,err := stateHandler.ExecuteChecked(util.PanicCheckedRunnable(&util.IntendedPanic{Message: uniqPanic}))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), uniqPanic)

	assert.NotEqual(t, stateHandler, container.getStateHandler())
	assert.False(
		t,
		reflect.TypeOf(container.getStateHandler()).AssignableTo(reflect.TypeOf(&CloseStateHandler{})),
	)
}

func TestCloseStateHandler_stillCloseAfterNumberOfAck(t *testing.T) {
	//init
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
		executed,err := stateHandler.ExecuteChecked(util.CheckedRunnable())
		assert.True(t, executed)
		assert.Nil(t, err)
	}

	assert.True(t, stateHandler.AcquirePermission())
	assert.Equal(t, stateHandler, container.getStateHandler())
	assert.True(
		t,
		reflect.TypeOf(container.getStateHandler()).AssignableTo(reflect.TypeOf(&CloseStateHandler{})),
	)

	errorAttempt := int(math.Ceil(float64(ctx.Config.SlidingWindowMaxSize)*float64(ctx.Config.ErrorThreshold)) - 1)
	for i:=0; i < errorAttempt; i++ {
		uniqPanic := util.RandErrorWithMessage()
		executed,result,err := stateHandler.ExecuteCheckedSupplier(util.ErrorCheckedSupplier(&uniqPanic))
		assert.True(t, executed)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), uniqPanic.Message)
	}

	assert.Equal(t, stateHandler, container.getStateHandler())
}