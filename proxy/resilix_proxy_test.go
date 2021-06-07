package proxy

import (
	"github.com/alfian853/resilix-go/context"
	"github.com/alfian853/resilix-go/statehandler"
	"github.com/alfian853/resilix-go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var proxyTestPanic = testutil.RandErrorWithMessage()

type stateHandlerMock struct {
	testutil.ModifiedMock
	statehandler.StateHandler
}

func (mock *stateHandlerMock) EvaluateState() {
	mock.Called()
}

func (mock *stateHandlerMock) ExecuteChecked(fun func() error) (executed bool, err error) {
	args := mock.Called(fun)
	return args.Bool(0), args.Error(1)
}

func (mock *stateHandlerMock) ExecuteCheckedSupplier(fun func() (interface{}, error)) (executed bool, result interface{}, err error) {
	args := mock.Called(fun)
	return args.Bool(0), args.Get(1), args.Error(2)
}

func getMockedStateHandler(isRuntimePanic bool) (positiveMock *stateHandlerMock, negativeMock *stateHandlerMock) {
	positiveMock = new(stateHandlerMock)
	negativeMock = new(stateHandlerMock)

	positiveMock.On("EvaluateState")
	positiveMock.On("ExecuteChecked", mock.Anything).Return(true, nil)
	positiveMock.On("ExecuteCheckedSupplier", mock.Anything).Return(true, true, nil)

	if isRuntimePanic {
		negativeMock.On("EvaluateState")
		negativeMock.On("ExecuteChecked", mock.Anything).Panic(proxyTestPanic.Message)
		negativeMock.On("ExecuteCheckedSupplier", mock.Anything).Panic(proxyTestPanic.Message)
	} else {
		negativeMock.On("EvaluateState")
		negativeMock.On("ExecuteChecked", mock.Anything).Return(true, proxyTestPanic)
		negativeMock.On("ExecuteCheckedSupplier", mock.Anything).Return(true, nil, proxyTestPanic)
	}

	return positiveMock, negativeMock
}

func TestResilixProxy_executeChecked(t *testing.T) {

	ctx := context.NewContextDefault()
	proxy := NewResilixProxy(ctx)

	positiveStateHandler, negativeStateHandler := getMockedStateHandler(false)

	proxy.SetStateHandler(positiveStateHandler)

	executed, err := proxy.ExecuteChecked(testutil.CheckedRunnable())
	assert.True(t, executed)
	assert.Nil(t, err)

	executed, result, err := proxy.ExecuteCheckedSupplier(testutil.TrueCheckedSupplier())
	assert.True(t, executed)
	assert.True(t, result.(bool))
	assert.Nil(t, err)

	assert.Same(t, positiveStateHandler, proxy.GetStateHandler())

	positiveStateHandler.AssertNumberOfCalls(t, "EvaluateState", 3)
	positiveStateHandler.AssertCalled(t, "ExecuteChecked", mock.Anything)
	positiveStateHandler.AssertCalled(t, "ExecuteCheckedSupplier", mock.Anything)

	positiveStateHandler.On("EvaluateState").Run(func(args mock.Arguments) {
		proxy.SetStateHandler(negativeStateHandler)
	})

	// will trigger stateHandler.EvaluateState()
	assert.Same(t, negativeStateHandler, proxy.GetStateHandler())

	executed, result, err = proxy.ExecuteCheckedSupplier(testutil.ErrorCheckedSupplier(proxyTestPanic))

	assert.True(t, executed)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), proxyTestPanic.Error())

	assert.Same(t, negativeStateHandler, proxy.GetStateHandler())

	negativeStateHandler.AssertCalled(t, "ExecuteCheckedSupplier", mock.Anything)
	negativeStateHandler.AssertNumberOfCalls(t, "EvaluateState", 2)
}

func TestResilixProxy_executeUnsafe(t *testing.T) {

	ctx := context.NewContextDefault()
	proxy := NewResilixProxy(ctx)

	positiveStateHandler, negativeStateHandler := getMockedStateHandler(true)

	proxy.SetStateHandler(positiveStateHandler)

	executed := proxy.Execute(testutil.DoNothingRunnable())
	assert.True(t, executed)

	executed, result := proxy.ExecuteSupplier(testutil.TrueSupplier())
	assert.True(t, executed)
	assert.True(t, result.(bool))

	assert.Same(t, positiveStateHandler, proxy.GetStateHandler())

	positiveStateHandler.AssertNumberOfCalls(t, "EvaluateState", 3)
	positiveStateHandler.AssertCalled(t, "ExecuteChecked", mock.Anything)
	positiveStateHandler.AssertCalled(t, "ExecuteCheckedSupplier", mock.Anything)

	positiveStateHandler.On("EvaluateState").Run(func(args mock.Arguments) {
		proxy.SetStateHandler(negativeStateHandler)
	})

	// will trigger stateHandler.EvaluateState()
	assert.Same(t, negativeStateHandler, proxy.GetStateHandler())

	assert.PanicsWithValue(t, proxyTestPanic.Message, func() {
		proxy.ExecuteCheckedSupplier(testutil.ErrorCheckedSupplier(proxyTestPanic))
	})

	assert.Same(t, negativeStateHandler, proxy.GetStateHandler())

	negativeStateHandler.AssertCalled(t, "ExecuteCheckedSupplier", mock.Anything)
	negativeStateHandler.AssertNumberOfCalls(t, "EvaluateState", 2)
}
