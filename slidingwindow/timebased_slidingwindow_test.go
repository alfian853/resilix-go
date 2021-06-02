package slidingwindow

import (
	"github.com/stretchr/testify/assert"
	conf "resilix-go/config"
	"resilix-go/util"
	"sync"
	"testing"
	"time"
)

const(
	windowTimeRange = 250
)

func initTimeBasedSlidingWindow() *TimeBasedSlidingWindow {
	config := conf.NewConfiguration()
	config.SlidingWindowTimeRange = windowTimeRange
	return NewTimeBasedSlidingWindow(config)
}

//testcase: fire with 25 random ack followed by 10(70% success) ack in arbitrary order
func TestCompleteCase(t *testing.T){
	twindow := initTimeBasedSlidingWindow()
	var wg sync.WaitGroup
	assert.Equal(t, float32(0.0), twindow.GetErrorRate())

	for i := 0 ; i < 25; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			twindow.AckAttempt(util.RandBool())
		}, &wg)
	}
	time.Sleep(windowTimeRange * time.Millisecond)

	nSuccess := 7
	nFailure := 3

	for nSuccess > 0 || nFailure > 0 {
		if util.RandBool() || nFailure == 0 {
			nSuccess--
			wg.Add(1)
			util.AsyncWgRunner(func() {
				twindow.AckAttempt(true)
			}, &wg)
		} else {
			nFailure--
			wg.Add(1)
			util.AsyncWgRunner(func() {
				twindow.AckAttempt(false)
			}, &wg)
		}
	}
	wg.Wait()
	assert.InEpsilon(t, 0.3, twindow.GetErrorRate(), EPSILON)
	time.Sleep((windowTimeRange + 10) * time.Millisecond)
	assert.Equal(t, 0, twindow.getQueSize())
}

func TestTimeBasedSwObserver(t *testing.T){
	twindow := initTimeBasedSlidingWindow()
	var wg sync.WaitGroup

	count1 := util.NewInt32(0)
	count2 := util.NewInt32(0)
	count3 := util.NewInt32(0)

	observer1 := mockObserver{name: "obs-1",count: count1}
	observer2 := mockObserver{name: "obs-2",count: count2}
	observer3 := mockObserver{name: "obs-3",count: count3}

	twindow.AddObserver(observer1)
	twindow.AddObserver(observer2)
	twindow.AddObserver(observer3)

	for i := 0 ; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			twindow.AckAttempt(util.RandBool())
		}, &wg)
	}
	wg.Wait()
	twindow.RemoveObserver(observer1)

	for i := 0 ; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			twindow.AckAttempt(util.RandBool())
		}, &wg)
	}

	wg.Wait()
	twindow.RemoveObserver(observer2)

	for i := 0 ; i < 5; i++ {
		wg.Add(1)
		util.AsyncWgRunner(func() {
			twindow.AckAttempt(util.RandBool())
		}, &wg)
	}
	wg.Wait()
	assert.Equal(t, int32(5), *count1)
	assert.Equal(t, int32(10), *count2)
	assert.Equal(t, int32(15), *count3)
	assert.Equal(t, 15, twindow.getQueSize())
	time.Sleep((windowTimeRange + 10) * time.Millisecond)
	assert.Equal(t, 0, twindow.getQueSize())
}

