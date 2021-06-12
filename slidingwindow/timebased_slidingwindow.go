package slidingwindow

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/util"
	"github.com/oleiade/lane"
	"sync/atomic"
)

type TimeBasedSlidingWindow struct {
	DefaultSlidingWindow
	lock       SwLock
	cfg        *config.Configuration
	successQue *lane.Deque
	failureQue *lane.Deque
}

func NewTimeBasedSlidingWindow(cfg *config.Configuration) *TimeBasedSlidingWindow {
	swindow := TimeBasedSlidingWindow{}
	swindow.lock = util.NewInt32(SwLock_Available)
	swindow.cfg = cfg
	swindow.successQue = lane.NewDeque()
	swindow.failureQue = lane.NewDeque()
	swindow.DefaultSlidingWindow.Decorate(&swindow, cfg)
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
	return float32(failureSize) / float32(successSize+failureSize)
}

func (window *TimeBasedSlidingWindow) Clear() {
	defer atomic.SwapInt32(window.lock, SwLock_Available)

	if atomic.SwapInt32(window.lock, SwLock_ClearingAll) < SwLock_ClearingAll {
		window.successQue = lane.NewDeque()
		window.failureQue = lane.NewDeque()
	}
}

func (window *TimeBasedSlidingWindow) examineAttemptWindow() {

	defer atomic.SwapInt32(window.lock, SwLock_Available)

	if atomic.SwapInt32(window.lock, SwLock_Clearing) < SwLock_Clearing {
		for !window.successQue.Empty() &&
			window.successQue.First().(int64) < util.GetTimestamp()-window.cfg.SlidingWindowTimeRange {
			window.successQue.Shift()
		}
		for !window.failureQue.Empty() &&
			window.failureQue.First().(int64) < util.GetTimestamp()-window.cfg.SlidingWindowTimeRange {
			window.failureQue.Shift()
		}

	}
}
