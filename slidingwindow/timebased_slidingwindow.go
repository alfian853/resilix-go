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
}

func NewTimeBasedSlidingWindow(config *conf.Configuration) *TimeBasedSlidingWindow {
	swindow := TimeBasedSlidingWindow{}
	swindow.lock = util.NewInt32(consts.SwLock_Available)
	swindow.config = config
	swindow.successQue = lane.NewDeque()
	swindow.failureQue = lane.NewDeque()
	swindow.DefaultSlidingWindow.Decorate(&swindow, config)
	return &swindow
}

func (window *TimeBasedSlidingWindow) handleAckAttempt(success bool) {
	lastAttempt := util.GetTimestamp()

	if success {
		window.successQue.Append(lastAttempt)
	} else {
		window.failureQue.Append(lastAttempt)
	}

	window.examineAttemptWindow()
}

func (window *TimeBasedSlidingWindow) getQueSize() int {
	window.examineAttemptWindow()
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
			window.successQue.First().(int64) < util.GetTimestamp() -window.config.SlidingWindowTimeRange {
			window.successQue.Shift()
		}
		for !window.failureQue.Empty() &&
			window.failureQue.First().(int64) < util.GetTimestamp() -window.config.SlidingWindowTimeRange {
			window.failureQue.Shift()
		}

	}
}