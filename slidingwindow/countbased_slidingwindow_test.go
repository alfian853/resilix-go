package slidingwindow

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	windowSize = 10
)

func initCountBasedSlidingWindow() *CountBasedSlidingWindow {
	config := conf.NewConfiguration()
	config.SlidingWindowMaxSize = windowSize
	return NewCountBasedSlidingWindow(config)
}

//testcase: fire with 25 random ack followed by 10(70% success) ack in arbitrary order
func TestCountBasedSwCompleteCase(t *testing.T) {
	cwindow := initCountBasedSlidingWindow()
	var wg sync.WaitGroup
	assert.Equal(t, float32(0.0), cwindow.GetErrorRate())

	for i := 0; i < 25; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			cwindow.AckAttempt(testutil.RandBool())
		}, &wg)
	}
	wg.Wait()
	nSuccess := 7
	nFailure := 3

	for nSuccess > 0 || nFailure > 0 {
		if testutil.RandBool() || nFailure == 0 {
			nSuccess--
			wg.Add(1)
			util.AsyncWgRunner(func() {
				cwindow.AckAttempt(true)
			}, &wg)
		} else {
			nFailure--
			wg.Add(1)
			util.AsyncWgRunner(func() {
				cwindow.AckAttempt(false)
			}, &wg)
		}
	}
	wg.Wait()
	assert.InEpsilon(t, 0.3, cwindow.GetErrorRate(), EPSILON)
}

func (obs mockObserver) NotifyOnAckAttempt(success bool) {
	atomic.AddInt32(obs.count, 1)
}

func TestCountBasedSwObserver(t *testing.T) {
	cwindow := initCountBasedSlidingWindow()
	var wg sync.WaitGroup

	count1 := util.NewInt32(0)
	count2 := util.NewInt32(0)
	count3 := util.NewInt32(0)

	observer1 := mockObserver{name: "obs-1", count: count1}
	observer2 := mockObserver{name: "obs-2", count: count2}
	observer3 := mockObserver{name: "obs-3", count: count3}

	cwindow.AddObserver(observer1)
	cwindow.AddObserver(observer2)
	cwindow.AddObserver(observer3)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			cwindow.AckAttempt(testutil.RandBool())
		}, &wg)
	}
	wg.Wait()
	cwindow.RemoveObserver(observer1)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			cwindow.AckAttempt(testutil.RandBool())
		}, &wg)
	}

	wg.Wait()
	cwindow.RemoveObserver(observer2)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			cwindow.AckAttempt(testutil.RandBool())
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, int32(5), *count1)
	assert.Equal(t, int32(10), *count2)
	assert.Equal(t, int32(15), *count3)
}
