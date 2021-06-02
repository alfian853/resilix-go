package slidingwindow

import (
	"github.com/oleiade/lane"
	conf "resilix-go/config"
	"resilix-go/consts"
	"resilix-go/util"
	"sync/atomic"
)


type TimeBasedSlidingWindow struct {
	DefaultSlidingWindow
	lock        consts.SwLock
	config      *conf.Configuration
	successQue  *lane.Deque
	failureQue  *lane.Deque
	lastAttempt *int64
}

func NewTimeBasedSlidingWindow(config *conf.Configuration){
	swindow := TimeBasedSlidingWindow{}
	*swindow.lock = consts.SwLock_Available
	swindow.config = config
	swindow.successQue = &lane.Deque{}
	swindow.failureQue = &lane.Deque{}
	*swindow.lastAttempt = 0
}

func (window *TimeBasedSlidingWindow) handleAckAttempt(success bool) {
	lastAttempt := util.GetTimestamp()
	atomic.SwapInt64(window.lastAttempt, lastAttempt)

	if success {
		window.successQue.Append(window.lastAttempt)
	} else {
		window.failureQue.Append(window.lastAttempt)
	}

	window.examineAttemptWindow()
}

func (window *TimeBasedSlidingWindow) getQueSize() int {
	return window.successQue.Size() + window.failureQue.Size()
}

func (window *TimeBasedSlidingWindow) getErrorRateAfterMinCallSatisfied() float32 {

	if window.successQue.Empty() && window.failureQue.Empty() {
		return 0
	}

	successSize, failureSize :=
		window.successQue.Size(), window.failureQue.Size()
	return float32(failureSize) / float32(successSize +failureSize)
}

func (window *TimeBasedSlidingWindow) Clear() {
	defer atomic.SwapInt32(window.lock, consts.SwLock_Available)

	if atomic.SwapInt32(window.lock, consts.SwLock_ClearingAll) < consts.SwLock_ClearingAll {
		window.successQue = &lane.Deque{}
		window.failureQue = &lane.Deque{}
	}
}

func (window *TimeBasedSlidingWindow) examineAttemptWindow() {

	defer atomic.SwapInt32(window.lock, consts.SwLock_Available)

	if atomic.SwapInt32(window.lock, consts.SwLock_Clearing) < consts.SwLock_Clearing {
		for !window.successQue.Empty() &&
			window.successQue.First().(int64) < (atomic.LoadInt64(window.lastAttempt) -window.config.SlidingWindowTimeRange) {
			window.successQue.Shift()
		}
		for !window.failureQue.Empty() &&
			window.failureQue.First().(int64) < (atomic.LoadInt64(window.lastAttempt) -window.config.SlidingWindowTimeRange) {
			window.failureQue.Shift()
		}

	}
}