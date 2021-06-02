package slidingwindow

import (
	"github.com/oleiade/lane"
	conf "resilix-go/config"
	"resilix-go/consts"
	"resilix-go/util"
	"sync/atomic"
)

type CountBasedSlidingWindow struct {
	DefaultSlidingWindow
	DefaultSlidingWindowExt
	lock       consts.SwLock
	config     *conf.Configuration
	errorCount *int32
	windowQue  *lane.Deque
}

func NewCountBasedSlidingWindow(config *conf.Configuration) *CountBasedSlidingWindow {
	swindow := CountBasedSlidingWindow{}
	swindow.lock = util.NewInt32(consts.SwLock_Available)
	swindow.errorCount = util.NewInt32(0)
	swindow.config = config
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
	defer atomic.SwapInt32(swindow.lock, consts.SwLock_Available)

	if atomic.SwapInt32(swindow.lock, consts.SwLock_ClearingAll) < consts.SwLock_ClearingAll {
		swindow.windowQue = &lane.Deque{}
	}
}

func (swindow *CountBasedSlidingWindow) examineAttemptWindow() {

	defer atomic.SwapInt32(swindow.lock, consts.SwLock_Available)

	if atomic.SwapInt32(swindow.lock, consts.SwLock_Clearing) < consts.SwLock_Clearing {
		for swindow.windowQue.Size() > swindow.config.SlidingWindowMaxSize {
			if swindow.windowQue.Shift() == false {
				atomic.AddInt32(swindow.errorCount, -1)
			}
		}
	}
}