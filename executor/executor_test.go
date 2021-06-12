package executor

import (
	"github.com/alfian853/resilix-go/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockDefaultExecutorExt struct {
	DefaultExecutorExt
	permission   bool
	successCount int
	failureCount int
}

func (ext *mockDefaultExecutorExt) AcquirePermission() bool {
	return ext.permission
}

func (ext *mockDefaultExecutorExt) OnAfterExecution(success bool) {
	if success {
		ext.successCount++
	} else {
		ext.failureCount++
	}
}

func TestDefaultExecutor_ExecuteCheckedSupplier(t *testing.T) {
	executorExt := mockDefaultExecutorExt{permission: true}
	executor := new(DefaultExecutor).Decorate(&executorExt)

	expectedResult := "expected-result"
	expectedError := testutil.RandErrorWithMessage()

	executed, result, err := executor.ExecuteCheckedSupplier(func() (interface{}, error) {
		return expectedResult, expectedError
	})

	assert.True(t, executed)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedError.Error(), err.Error())

	// set permission = false
	executorExt.permission = false

	executed, result, err = executor.ExecuteCheckedSupplier(func() (interface{}, error) {
		return expectedResult, expectedError
	})
	assert.False(t, executed)
	assert.Nil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, 0, executorExt.successCount)
	assert.Equal(t, 1, executorExt.failureCount)
}

func TestDefaultExecutor_ExecuteChecked(t *testing.T) {

	executorExt := mockDefaultExecutorExt{permission: true}
	executor := new(DefaultExecutor).Decorate(&executorExt)

	expectedError := testutil.RandErrorWithMessage()
	executed, err := executor.ExecuteChecked(testutil.PanicCheckedRunnable(expectedError))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), expectedError.Error())

	expectedError = testutil.RandErrorWithMessage()
	executed, err = executor.ExecuteChecked(testutil.ErrorCheckedRunnable(expectedError))
	assert.True(t, executed)
	assert.Same(t, expectedError, err)

	isRun := false
	executed, err = executor.ExecuteChecked(testutil.SpyCheckedRunnable(&isRun))
	assert.True(t, executed)
	assert.True(t, isRun)
	assert.Nil(t, err)

	isRun = false
	executorExt.permission = false
	executed, err = executor.ExecuteChecked(testutil.SpyCheckedRunnable(&isRun))
	assert.False(t, executed)
	assert.False(t, isRun)
	assert.Nil(t, err)

	assert.Equal(t, 1, executorExt.successCount)
	assert.Equal(t, 2, executorExt.failureCount)
}
