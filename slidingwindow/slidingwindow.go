package slidingwindow

import (
	conf "github.com/alfian853/resilix-go/config"
	"sync"
)

type SwObserver interface {
	NotifyOnAckAttempt(success bool)
}

type SlidingWindow interface {
	AckAttempt(success bool)
	GetErrorRate() float32
	Clear()

	AddObserver(observer SwObserver)
	RemoveObserver(observer SwObserver)
}

type DefaultSlidingWindowExt interface {
	handleAckAttempt(success bool)
	getQueSize() int
	getErrorRateAfterMinCallSatisfied() float32
}

type DefaultSlidingWindow struct {
	SlidingWindow
	swindowExt   DefaultSlidingWindowExt
	config       *conf.Configuration
	observerLock sync.Mutex
	observers    []SwObserver
}

func (swindow *DefaultSlidingWindow) Decorate(swindowExt DefaultSlidingWindowExt, config *conf.Configuration) {
	swindow.swindowExt = swindowExt
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

	swindow.swindowExt.handleAckAttempt(success)
	for i := 0; i < len(swindow.observers); i++ {
		swindow.observers[i].NotifyOnAckAttempt(success)
	}
}

func (swindow *DefaultSlidingWindow) GetErrorRate() float32 {
	if swindow.swindowExt.getQueSize() < swindow.config.MinimumCallToEvaluate {
		return 0
	}

	return swindow.swindowExt.getErrorRateAfterMinCallSatisfied()
}
