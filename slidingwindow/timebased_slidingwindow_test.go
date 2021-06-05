package slidingwindow

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/alfian853/resilix-go/util"
	"sync"
	"testing"
	"time"
)

const(
	windowTimeRange = 550
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
			twindow.AckAttempt(testutil.RandBool())
		}, &wg)
	}
	wg.Wait()
	time.Sleep((windowTimeRange+100) * time.Millisecond)

	nSuccess := 7
	nFailure := 3

	for nSuccess > 0 || nFailure > 0 {
		if testutil.RandBool() || nFailure == 0 {
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