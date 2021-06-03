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


func TestPessimisticRetryRejected(t *testing.T){
	t.Deadline()
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.3
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_PESSIMISTIC
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
	retryManager := NewPessimisticRetryManager().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryManager.GetRetryState()))
	assert.Equal(t, float32(0), retryManager.GetErrorRate())

	retryCount := 0

	for retryManager.GetErrorRate() < ctx.Config.ErrorThreshold {
		if retryManager.AcquireAndUpdateRetryPermission() {
			wg.Add(1)
			util.AsyncWgRunner(func() {
				retryManager.AcquireAndUpdateRetryPermission()
				window.AckAttempt(util.RandBool())
			}, &wg)
			retryCount++
		}
	}
	wg.Wait()

	assert.True(t, retryManager.GetErrorRate() >= ctx.Config.ErrorThreshold)
	assert.True(t, int(ctx.Config.NumberOfRetryInHalfOpenState) > retryCount)
}

func TestPessimisticRetryAcceptedCase(t *testing.T) {
	//init
	ctx := context.NewContextDefault()
	ctx.Config.SlidingWindowMaxSize = 50
	ctx.Config.ErrorThreshold = 0.8
	ctx.Config.NumberOfRetryInHalfOpenState = 100
	ctx.Config.RetryStrategy = consts.RETRY_PESSIMISTIC
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
	retryManager := NewPessimisticRetryManager().Decorate(ctx)

	assert.Equal(t, consts.RETRY_ON_GOING, int(retryManager.GetRetryState()))
	assert.Equal(t, float32(0), retryManager.GetErrorRate())

	minSuccessAck := int(((1 - ctx.Config.ErrorThreshold) * float32(ctx.Config.NumberOfRetryInHalfOpenState)) + 2)

	for i:=0 ; i < minSuccessAck && (retryManager.GetErrorRate() < ctx.Config.ErrorThreshold); i++ {

		if retryManager.AcquireAndUpdateRetryPermission() {
			wg.Add(1)
			util.AsyncWgRunner(func() {
				window.AckAttempt(true)
			}, &wg)

		} else {
			i--
		}

	}

	for i:=0 ; i < int(ctx.Config.NumberOfRetryInHalfOpenState) - minSuccessAck; i++ {
		if retryManager.AcquireAndUpdateRetryPermission() {
			wg.Add(1)
			util.AsyncWgRunner(func() {
				window.AckAttempt(false)
			}, &wg)

		} else {
			i--
		}
	}

	wg.Wait()

	assert.False(t,  retryManager.AcquireAndUpdateRetryPermission())
	assert.Equal(t, consts.RETRY_ACCEPTED, int(retryManager.GetRetryState()))
}