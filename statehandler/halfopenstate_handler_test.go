package statehandler

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/slidingwindow"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHalfOpenState_retryAndSuccess(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5 //be careful to edit this, see other comments below
	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 5 //be careful to edit this, see other comments below
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := new(HalfOpenStateHandler).Decorate(ctx, &container)

	container.SetStateHandler(stateHandler)

	// maxAcceptableError is 2 since the error threshold is 50% for 5 times retry
	maxAcceptableError := 2
	// shouldSuccessAttempt is 3
	shouldSuccessAttempt := 3

	for i := 0; i < shouldSuccessAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed, result, err := stateHandler.ExecuteCheckedSupplier(testutil.TrueCheckedSupplier())
			assert.True(t, executed)
			assert.True(t, result.(bool))
			assert.Nil(t, err)
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.GetStateHandler())

	for i := 0; i < maxAcceptableError; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqueError := testutil.RandPanicMessage()
			executed, err := stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable(uniqueError))
			assert.True(t, executed)
			assert.Contains(t, err.Error(), uniqueError)
		}, &wg)
	}
	wg.Wait()
	assert.NotEqual(t, stateHandler, container.GetStateHandler())
	assert.IsType(t, &CloseStateHandler{}, container.GetStateHandler())
	assert.Equal(t, float32(0), ctx.SWindow.GetErrorRate())
}

func TestHalfOpenState_retryAndFailed(t *testing.T) {
	//init
	var wg sync.WaitGroup
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5 //be careful to edit this, see other comments below
	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 6 //be careful to edit this, see other comments below
	ctx.Config.WaitDurationInOpenState = 100000000
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := new(HalfOpenStateHandler).Decorate(ctx, &container)

	container.SetStateHandler(stateHandler)

	// minRequiredError is 2 since the error threshold is 50% for 5 times retry
	minRequiredError := 3
	shouldSuccessAttempt := 3

	for i := 0; i < shouldSuccessAttempt; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			executed, result, err := stateHandler.ExecuteCheckedSupplier(testutil.TrueCheckedSupplier())
			assert.True(t, executed)
			assert.True(t, result.(bool))
			assert.Nil(t, err)
		}, &wg)

	}
	wg.Wait()
	assert.Equal(t, stateHandler, container.GetStateHandler())

	for i := 0; i < minRequiredError; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			uniqueError := testutil.RandPanicMessage()
			executed, err := stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable(uniqueError))
			assert.True(t, executed)
			assert.Contains(t, err.Error(), uniqueError)
		}, &wg)

	}
	wg.Wait()
	uniqueError := testutil.RandPanicMessage()
	executed, err := stateHandler.ExecuteChecked(testutil.PanicCheckedRunnable(uniqueError))
	assert.False(t, executed)
	assert.Nil(t, err)

	assert.NotEqual(t, stateHandler, container.GetStateHandler())
	assert.IsType(t, &OpenStateHandler{}, container.GetStateHandler())
}
