package slidingwindow

import (
	"github.com/stretchr/testify/assert"
	conf "resilix-go/config"
	"resilix-go/slidingwindow/test"
	"resilix-go/util"
	"sync/atomic"
	"testing"
)

const(
	WINDOW_SIZE = 10
)

func initSlidingWindow() *CountBasedSlidingWindow {
	config := conf.NewConfiguration()
	config.SlidingWindowMaxSize = WINDOW_SIZE
	return NewCountBasedSlidingWindow(config)
}

//testcase: fire with 25 random ack followed by 10(70% success) ack in arbitrary order
func TestCompleteCase(t *testing.T){
	cwindow := initSlidingWindow()
	assert.Equal(t, float32(0.0), cwindow.GetErrorRate())

	for i := 0 ; i < 25; i++ {
		cwindow.AckAttempt(util.RandBool())
	}

	nSuccess := 7
	nFailure := 3

	for nSuccess > 0 || nFailure > 0 {
		if util.RandBool() || nFailure == 0 {
			nSuccess--
			cwindow.AckAttempt(true)
		} else {
			nFailure--
			cwindow.AckAttempt(false)
		}
	}

	assert.InEpsilon(t, 0.3, cwindow.GetErrorRate(), test.EPSILON)
}

type mockObserver struct {
	SwObserver
	name string
	count *int32
}

func (obs mockObserver) NotifyOnAckAttempt(success bool) {
	atomic.AddInt32(obs.count, 1)
}

func TestObserver(t *testing.T){
	cwindow := initSlidingWindow()
	count1 := util.NewInt32(0)
	count2 := util.NewInt32(0)
	count3 := util.NewInt32(0)

	observer1 := mockObserver{name: "obs-1",count: count1}
	observer2 := mockObserver{name: "obs-2",count: count2}
	observer3 := mockObserver{name: "obs-3",count: count3}

	cwindow.AddObserver(observer1)
	cwindow.AddObserver(observer2)
	cwindow.AddObserver(observer3)

	for i := 0 ; i < 5; i++ {
		cwindow.AckAttempt(util.RandBool())
	}

	cwindow.RemoveObserver(observer1)

	for i := 0 ; i < 5; i++ {
		cwindow.AckAttempt(util.RandBool())
	}

	cwindow.RemoveObserver(observer2)

	for i := 0 ; i < 5; i++ {
		cwindow.AckAttempt(util.RandBool())
	}

	assert.Equal(t, int32(5), *count1)
	assert.Equal(t, int32(10), *count2)
	assert.Equal(t, int32(15), *count3)
}

