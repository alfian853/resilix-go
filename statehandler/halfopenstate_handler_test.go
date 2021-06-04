package statehandler

import (
	"github.com/stretchr/testify/assert"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/util"
	"sync"
	"testing"
)

func TestHalfOpenState_retryAndSuccess(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5 //be careful to edit this, see other comments below
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 5 //be careful to edit this, see other comments below
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewHalfOpenStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	// maxAcceptableError is 2 since the error threshold is 50% for 5 times retry
	maxAcceptableError := 2
	// shouldSuccessAttempt is 3
	shouldSuccessAttempt := 3

	for i:=0; i < shouldSuccessAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed,result,err := stateHandler.ExecuteCheckedSupplier(util.TrueCheckedSupplier())
			assert.True(t, executed)
			assert.True(t, result.(bool))
			assert.Nil(t, err)
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.getStateHandler())

	for i:=0; i < maxAcceptableError; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqueError := util.RandPanicMessage()
			executed,err := stateHandler.ExecuteChecked(util.PanicCheckedRunnable(uniqueError))
			assert.True(t, executed)
			assert.Contains(t, err.Error(), uniqueError)
		}, &wg)
	}
	wg.Wait()
	assert.NotEqual(t, stateHandler, container.getStateHandler())
	assert.True(t, container.getStateHandler().AcquirePermission())
	assert.IsType(t, &CloseStateHandler{}, container.getStateHandler())
	assert.Equal(t, float32(0), ctx.SWindow.GetErrorRate())
}

func TestHalfOpenState_retryAndFailed(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5 //be careful to edit this, see other comments below
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 6 //be careful to edit this, see other comments below
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := NewHalfOpenStateHandler().Decorate(ctx, &container)

	container.setStateHandler(stateHandler)

	// minRequiredError is 2 since the error threshold is 50% for 5 times retry
	minRequiredError := 3
	shouldSuccessAttempt := 3

	for i:=0; i < shouldSuccessAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed,result,err := stateHandler.ExecuteCheckedSupplier(util.TrueCheckedSupplier())
			assert.True(t, executed)
			assert.True(t, result.(bool))
			assert.Nil(t, err)
		}, &wg)

	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.getStateHandler())

	for i:=0; i < minRequiredError; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqueError := util.RandPanicMessage()
			executed,err := stateHandler.ExecuteChecked(util.PanicCheckedRunnable(uniqueError))
			assert.True(t, executed)
			assert.Contains(t, err.Error(), uniqueError)
		}, &wg)

	}
	wg.Wait()
	uniqueError := util.RandPanicMessage()
	executed,err := stateHandler.ExecuteChecked(util.PanicCheckedRunnable(uniqueError))
	assert.False(t, executed)
	assert.Nil(t, err)

	assert.NotEqual(t, stateHandler, container.getStateHandler())
	assert.IsType(t, &OpenStateHandler{}, container.getStateHandler())
	assert.False(t, container.getStateHandler().AcquirePermission())
}