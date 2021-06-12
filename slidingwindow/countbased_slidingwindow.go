package slidingwindow

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/util"
	"github.com/oleiade/lane"
	"sync/atomic"
)

type CountBasedSlidingWindow struct {
	DefaultSlidingWindow
	DefaultSlidingWindowExt
	lock       SwLock
	cfg        *config.Configuration
	errorCount *int32
	windowQue  *lane.Deque
}

func NewCountBasedSlidingWindow(config *config.Configuration) *CountBasedSlidingWindow {
	swindow := CountBasedSlidingWindow{}
	swindow.lock = util.NewInt32(SwLock_Available)
	swindow.errorCount = util.NewInt32(0)
	swindow.cfg = config
	swindow.windowQue = lane.NewDeque()
	swindow.DefaultSlidingWindow.Decorate(&swindow, config)
	return &swindow
}

func (swindow *CountBasedSlidingWindow) handleAckAttempt(success bool) {
	swindow.windowQue.Append(success)

	if !success {
		atomic.AddInt32(swindow.errorCount, 1)
	}

	swindow.examineAttemptWindow()
}

func (swindow *CountBasedSlidingWindow) getQueSize() int {
	return swindow.windowQue.Size()
}

func (swindow *CountBasedSlidingWindow) getErrorRateAfterMinCallSatisfied() float32 {

	if swindow.windowQue.Size() == 0 {
		return 0
	}

	return float32(atomic.LoadInt32(swindow.errorCount)) / float32(swindow.windowQue.Size())
}

func (swindow *CountBasedSlidingWindow) Clear() {
	defer atomic.SwapInt32(swindow.lock, SwLock_Available)

	if atomic.SwapInt32(swindow.lock, SwLock_ClearingAll) < SwLock_ClearingAll {
		swindow.windowQue = lane.NewDeque()
	}
}

func (swindow *CountBasedSlidingWindow) examineAttemptWindow() {

	defer atomic.SwapInt32(swindow.lock, SwLock_Available)

	if atomic.SwapInt32(swindow.lock, SwLock_Clearing) < SwLock_Clearing {
		for swindow.windowQue.Size() > swindow.cfg.SlidingWindowMaxSize {
			if swindow.windowQue.Shift() == false {
				atomic.AddInt32(swindow.errorCount, -1)
			}
		}
	}
}
