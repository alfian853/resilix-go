package proxy

import (
	"fmt"
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResilixProxyIntegration_executeChecked(t *testing.T) {

	ctx := context.NewContextDefault()
	ctx.Config.MinimumCallToEvaluate = 1000
	proxy := NewResilixProxy(ctx)

	isRun := false
	executed, err := proxy.ExecuteChecked(testutil.SpyCheckedRunnable(&isRun))
	assert.True(t, executed)
	assert.True(t, isRun)
	assert.Nil(t, err)

	expectedResult := fmt.Sprintf("expectedResult-%d", testutil.RandInt(1, 100))
	executed, result, err := proxy.ExecuteCheckedSupplier(testutil.CheckedSupplier(expectedResult))
	assert.True(t, executed)
	assert.Equal(t, expectedResult, result.(string))
	assert.Nil(t, err)

	expectedError := testutil.RandErrorWithMessage()
	executed, err = proxy.ExecuteChecked(testutil.PanicCheckedRunnable(expectedError))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), expectedError.Error())

	expectedError = testutil.RandErrorWithMessage()
	executed, result, err = proxy.ExecuteCheckedSupplier(testutil.PanicCheckedSupplier(expectedError))
	assert.True(t, executed)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedError.Error())

}

func TestResilixProxyIntegration_execute(t *testing.T) {

	ctx := context.NewContextDefault()
	ctx.Config.MinimumCallToEvaluate = 1000
	proxy := NewResilixProxy(ctx)

	isReallyExecuted := false
	executed, err := proxy.Execute(testutil.SpyRunnable(&isReallyExecuted))
	assert.True(t, executed)
	assert.True(t, isReallyExecuted)
	assert.Nil(t, err)

	expectedResult := fmt.Sprintf("expectedResult-%d", testutil.RandInt(1, 100))
	executed, result, err := proxy.ExecuteSupplier(testutil.Supplier(expectedResult))
	assert.True(t, executed)
	assert.Equal(t, expectedResult, result.(string))
	assert.Nil(t, err)

	expectedError := testutil.RandErrorWithMessage()
	executed, err = proxy.Execute(testutil.PanicRunnable(expectedError))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), expectedError.Error())

	expectedError = testutil.RandErrorWithMessage()

	executed, result, err = proxy.ExecuteSupplier(testutil.PanicSupplier(expectedError))
	assert.True(t, executed)
	assert.Contains(t, err.Error(), expectedError.Error())

}
