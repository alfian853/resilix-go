package registry

import (
	"github.com/stretchr/testify/assert"
	"resilix-go/config"
	"resilix-go/consts"
	"testing"
)

func TestCreateUnknown(t *testing.T) {
	executor := GetResilixExecutor("random-1")
	assert.NotNil(t, executor)
	assert.Same(t, executor, GetResilixExecutor("random-1"))
}

func TestCreate(t *testing.T) {
	executor1 := Register("context-1", nil)
	assert.NotNil(t, executor1)
	assert.Same(t, executor1, GetResilixExecutor("context-1"))

	config2 := config.NewConfiguration()
	config2.SlidingWindowStrategy = consts.SwStrategy_TimeBased
	executor2 := Register("context-2", config2)

	assert.NotNil(t, executor2)
	assert.Same(t, executor2, GetResilixExecutor("context-2"))
}