package slidingwindow

import (
	"resilix-go/context"
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

type defaultSlidingWindowExt interface {
	handleAckAttempt(success bool)
	getQueSize() int
	getErrorRateAfterMinCallSatisfied() float32
}

type defaultSlidingWindow struct {
	SlidingWindow
	defaultSlidingWindowExt
	isActive IsActive
	config *context.Configuration
	observerLock sync.Mutex
	observers []SwObserver
}

func (sw *defaultSlidingWindow) AddObserver(observer SwObserver) {
	defer sw.observerLock.Unlock()
	sw.observerLock.Lock()

	sw.observers = append(sw.observers, observer)
}

func (sw *defaultSlidingWindow) RemoveObserver(observer SwObserver) {
	defer sw.observerLock.Unlock()
	sw.observerLock.Lock()

	var targetIndex *int

	for i := 0; i < len(sw.observers); i++ {
		if sw.observers[i] == observer {
			targetIndex = &i
			break
		}
	}

	if targetIndex != nil {
		lastIndex := len(sw.observers) - 1

		// swap target with the last index
		sw.observers[*targetIndex], sw.observers[lastIndex] =
			sw.observers[lastIndex - 1], sw.observers[*targetIndex]

		sw.observers = sw.observers[:lastIndex]
	}
}


func (sw *defaultSlidingWindow) AckAttempt(success bool) {

	if atomic.LoadInt32(sw.isActive) == Active {
		sw.handleAckAttempt(success)
	}
}

func (sw *defaultSlidingWindow) GetErrorRate() float32 {
	if sw.getQueSize() < sw.config.MinimumCallToEvaluate {
		return 0
	}

	return sw.getErrorRateAfterMinCallSatisfied()
}

