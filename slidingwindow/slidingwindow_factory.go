package slidingwindow

import (
	conf "github.com/alfian853/resilix-go/config"
	"github.com/alfian853/resilix-go/consts"
)

func CreateSlidingWindow(conf *conf.Configuration) SlidingWindow {
	switch conf.SlidingWindowStrategy {
	case consts.SwStrategy_CountBased:
		return NewCountBasedSlidingWindow(conf)
	case consts.SwStrategy_TimeBased:
		return NewTimeBasedSlidingWindow(conf)
	}

	panic("Unknown SlidingWindowStrategy: " + conf.SlidingWindowStrategy)
}
