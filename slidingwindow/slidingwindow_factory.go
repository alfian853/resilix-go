package slidingwindow

import (
	"github.com/alfian853/resilix-go/config"
)

func CreateSlidingWindow(cfg *config.Configuration) SlidingWindow {
	switch cfg.SlidingWindowStrategy {
	case config.SwStrategy_CountBased:
		return NewCountBasedSlidingWindow(cfg)
	case config.SwStrategy_TimeBased:
		return NewTimeBasedSlidingWindow(cfg)
	}

	panic("Unknown SlidingWindowStrategy: " + string(cfg.SlidingWindowStrategy))
}
