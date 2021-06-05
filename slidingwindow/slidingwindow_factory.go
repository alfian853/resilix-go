package slidingwindow

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/consts"
)

func CreateSlidingWindow(conf *conf.Configuration) SlidingWindow {
	switch conf.RetryStrategy {
	case consts.RETRY_OPTIMISTIC:
		return NewCountBasedSlidingWindow(conf)
	case consts.RETRY_PESSIMISTIC:
		return NewTimeBasedSlidingWindow(conf)
	}

	return nil
}

