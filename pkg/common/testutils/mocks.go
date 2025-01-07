package testutils

import "github.com/stretchr/testify/mock"

// MockContext is used to mock any context in requests as they are not needed for testing purposes
const MockContext = mock.Anything

// MockCalls is a wrapper around []*mock.Call
type MockCalls []*mock.Call

// Times sets how many times we expect each call to execute
func (mc MockCalls) Times(t int) MockCalls {
	for _, c := range mc {
		c.Times(t)
	}
	return mc
}

// Once expects calls to be called only one time
func (mc MockCalls) Once() MockCalls {
	return mc.Times(1)
}

// ReturnErr sets the given error as a last return parameter of the call with the given method
func (mc MockCalls) ReturnErr(method string, err error) MockCalls {
	for _, c := range mc {
		if c.Method == method {
			last := len(c.ReturnArguments) - 1
			c.ReturnArguments[last] = err
		}
	}

	return mc
}
