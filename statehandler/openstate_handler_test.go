package statehandler

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/slidingwindow"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOpenState_movingStateAfterWaitingDurationIsSatisfied(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 10
	ctx.Config.ErrorThreshold = 0.5
	ctx.Config.RetryStrategy = config.RetryStrategy_Optimistic
	ctx.Config.MinimumCallToEvaluate = 3
	ctx.Config.NumberOfRetryInHalfOpenState = 10
	ctx.Config.WaitDurationInOpenState = 200
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	container := testStateContainer{}

	stateHandler := new(OpenStateHandler).Decorate(ctx, &container)

	container.SetStateHandler(stateHandler)

	assert.Equal(t, stateHandler, container.GetStateHandler())

	time.Sleep(time.Duration(ctx.Config.WaitDurationInOpenState/2) * time.Millisecond)

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
