package slidingwindow

import (
	"github.com/alfian853/resilix-go/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlidingWindowFactory(t *testing.T) {
	cfg := config.NewConfiguration()

	cfg.SlidingWindowStrategy = config.SwStrategy_CountBased
	assert.IsType(t, &CountBasedSlidingWindow{}, CreateSlidingWindow(cfg))

	cfg.SlidingWindowStrategy = config.SwStrategy_TimeBased
	assert.IsType(t, &TimeBasedSlidingWindow{}, CreateSlidingWindow(cfg))

	cfg.SlidingWindowStrategy = "nothing"
	assert.PanicsWithValue(t, "Unknown SlidingWindowStrategy: nothing", func() {
		CreateSlidingWindow(cfg)
	})

}
