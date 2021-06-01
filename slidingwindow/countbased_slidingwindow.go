package slidingwindow

import (
	"github.com/oleiade/lane"
	"resilix-go/context"
	"sync/atomic"
)

const (
	Available = 1
	NotAvailable = 0
)

type CountBasedSlidingWindow struct {
	defaultSlidingWindow
	lock *int32
	config *context.Configuration
	errorCount *int32
	windowQue *lane.Deque
}

func (window *CountBasedSlidingWindow) handleAckAttempt(success bool) {
	window.windowQue.Append(success)

	if !success {
		atomic.AddInt32(window.errorCount, 1)
	}

	window.examineAttemptWindow()
}

func (window *CountBasedSlidingWindow) getQueSize() int {
	return window.windowQue.Size()
}

func (window *CountBasedSlidingWindow) getErrorRateAfterMinCallSatisfied() float32 {

	if window.windowQue.Size() == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(window.errorCount)) / float32(window.windowQue.Size())
}

func (window *CountBasedSlidingWindow) clear() {
	defer atomic.SwapInt32(window.lock, Available)

	if atomic.SwapInt32(window.lock, NotAvailable) == Available {
		window.windowQue.Pop()
	}
}

func (window *CountBasedSlidingWindow) examineAttemptWindow() {

	defer atomic.SwapInt32(window.lock, Available)

	if atomic.SwapInt32(window.lock, NotAvailable) == Available {
		for window.windowQue.Size() > window.config.SlidingWindowMaxSize {
			if window.windowQue.First() == false {
				atomic.AddInt32(window.errorCount, -1)
			}
		}
	}
}