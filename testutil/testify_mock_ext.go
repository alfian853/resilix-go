package testutil

import (
	"github.com/stretchr/testify/mock"
	"reflect"
)

type ModifiedMock struct {
	mock.Mock
}

func (m *ModifiedMock) On(methodName string, arguments ...interface{}) *mock.Call {

	lenObs := len(m.ExpectedCalls)
	var targetIndex *int

	// search and replace existing method mock
	for i := 0; i < lenObs; i++ {
		if methodName == m.ExpectedCalls[i].Method && isArgsEqual(m.ExpectedCalls[i].Arguments, arguments){
			targetIndex = &i
			break
		}
	}

	if targetIndex != nil {
		if lenObs == 1 {
			m.ExpectedCalls = m.ExpectedCalls[:0]
		}
		lastIndex := lenObs - 1

		// swap target with the last index
		m.ExpectedCalls[*targetIndex], m.ExpectedCalls[lastIndex] =
			m.ExpectedCalls[lastIndex], m.ExpectedCalls[*targetIndex]

		m.ExpectedCalls = m.ExpectedCalls[:lastIndex]
	}
	//search and replace done

	return m.Mock.On(methodName, arguments...)
}

// isArgsEqual compares arguments
func isArgsEqual(expected mock.Arguments, args []interface{}) bool {
	if len(expected) != len(args) {
		return false
	}
	for i, v := range args {
		if !reflect.DeepEqual(expected[i], v) {
			return false
		}
	}
	return true
}

