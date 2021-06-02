package slidingwindow

import (
	conf "resilix-go/config"
	"resilix-go/util"
	"sync"
	"sync/atomic"
)

type SwObserver interface {
	NotifyOnAckAttempt(success bool)
}

type SlidingWindow interface {
	AckAttempt(success bool)
	GetErrorRate() float32
	SetActive(isActive bool)
	Clear()

	AddObserver(observer SwObserver)
	RemoveObserver(observer SwObserver)
}

type IsActive *int32

const (
	Active   = 1
	Inactive = 0
)

type DefaultSlidingWindowExt interface {
	handleAckAttempt(success bool)
	getQueSize() int
	getErrorRateAfterMinCallSatisfied() float32
}

type DefaultSlidingWindow struct {
	SlidingWindow
	swindowExt DefaultSlidingWindowExt
	isActive IsActive
	config *conf.Configuration
	observerLock sync.Mutex
	observers []SwObserver
}

func (swindow *DefaultSlidingWindow) Decorate(swindowExt DefaultSlidingWindowExt,config *conf.Configuration)  {
	swindow.swindowExt = swindowExt
	swindow.isActive = util.NewInt32(Active)
	swindow.config = config
	swindow.observerLock = sync.Mutex{}
	swindow.observers = make([]SwObserver, 0)
}

func (swindow *DefaultSlidingWindow) AddObserver(observer SwObserver) {
	defer swindow.observerLock.Unlock()
	swindow.observerLock.Lock()

	swindow.observers = append(swindow.observers, observer)
}

func (swindow *DefaultSlidingWindow) RemoveObserver(observer SwObserver) {
	defer swindow.observerLock.Unlock()
	swindow.observerLock.Lock()

	lenObs := len(swindow.observers)
	var targetIndex *int

	for i := 0; i < lenObs; i++ {
		if swindow.observers[i] == observer {
			targetIndex = &i
			break
		}
	}

	if targetIndex != nil {
		if lenObs == 1 {
			swindow.observers = swindow.observers[:0]
			return
		}
		lastIndex := lenObs - 1

		// swap target with the last index
		swindow.observers[*targetIndex], swindow.observers[lastIndex] =
			swindow.observers[lastIndex], swindow.observers[*targetIndex]

		swindow.observers = swindow.observers[:lastIndex]
	}
}


func (swindow *DefaultSlidingWindow) AckAttempt(success bool) {

	if atomic.LoadInt32(swindow.isActive) == Active {
		swindow.swindowExt.handleAckAttempt(success)
		for i:=0; i < len(swindow.observers); i++ {
			swindow.observers[i].NotifyOnAckAttempt(success)
		}
	}
}

func (swindow *DefaultSlidingWindow) GetErrorRate() float32 {
	if swindow.swindowExt.getQueSize() < swindow.config.MinimumCallToEvaluate {
		return 0
	}

	return swindow.swindowExt.getErrorRateAfterMinCallSatisfied()
}

