package retry

import (
	"github.com/stretchr/testify/assert"
	"resilix-go/consts"
	"resilix-go/context"
	"resilix-go/slidingwindow"
	"resilix-go/util"
	"sync"
	"testing"
)


func TestRejected(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.3
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	var wg sync.WaitGroup
	window := ctx.SWindow

	// start
	for i:=0 ; i < ctx.Config.SlidingWindowMaxSize; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			window.AckAttempt(util.RandBool())
		}, &wg)
	}
	wg.Wait()
	retryManager := NewOptimisticRetryManager().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryManager.GetRetryState()))
	assert.Equal(t, float32(0), retryManager.GetErrorRate())

	minFailAck := int((ctx.Config.ErrorThreshold * float32(ctx.Config.NumberOfRetryInHalfOpenState)) + 2)
	maxSuccessAck := int(ctx.Config.NumberOfRetryInHalfOpenState) - minFailAck

	for i:=0 ; i < maxSuccessAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			assert.Equal(t, true, retryManager.AcquireAndUpdateRetryPermission())
			window.AckAttempt(true)
		}, &wg)
	}

	for i:=0 ; i < minFailAck; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			retryManager.AcquireAndUpdateRetryPermission()
			window.AckAttempt(false)
		}, &wg)
	}

	wg.Wait()

	assert.Equal(t, true, retryManager.GetErrorRate() >= ctx.Config.ErrorThreshold)
	assert.Equal(t, false, retryManager.AcquireAndUpdateRetryPermission())
	assert.Equal(t, consts.RETRY_REJECTED, int(retryManager.GetRetryState()))
}

func TestAcceptedCase(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.8
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_OPTIMISTIC
	ctx.SWindow = slidingwindow.NewCountBasedSlidingWindow(ctx.Config)
	var wg sync.WaitGroup
	window := ctx.SWindow

	// start
	for i:=0 ; i < ctx.Config.SlidingWindowMaxSize; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			window.AckAttempt(util.RandBool())
		}, &wg)
	}
	wg.Wait()
	retryManager := NewOptimisticRetryManager().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryManager.GetRetryState()))
	assert.Equal(t, float32(0), retryManager.GetErrorRate())

	minSuccessAck := int(((1 - ctx.Config.ErrorThreshold) * float32(ctx.Config.NumberOfRetryInHalfOpenState)) + 2)

	for i:=0 ; i < minSuccessAck; i++ {
		assert.Equal(t, true, retryManager.AcquireAndUpdateRetryPermission())
		wg.Add(1)
		util.AsyncWgRunner(func() {
			window.AckAttempt(true)
		}, &wg)
	}

	for i:=0 ; i < int(ctx.Config.NumberOfRetryInHalfOpenState) - minSuccessAck; i++ {
		assert.Equal(t, true, retryManager.AcquireAndUpdateRetryPermission())
		wg.Add(1)
		util.AsyncWgRunner(func() {
			window.AckAttempt(false)
		}, &wg)
	}

	wg.Wait()

	assert.Equal(t, true, retryManager.GetErrorRate() < ctx.Config.ErrorThreshold)
	assert.Equal(t, false, retryManager.AcquireAndUpdateRetryPermission())
	assert.Equal(t, consts.RETRY_ACCEPTED, int(retryManager.GetRetryState()))
}