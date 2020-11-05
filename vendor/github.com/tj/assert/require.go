/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package assert

import (
	http "net/http"
	url "net/url"
	time "time"

	assert "github.com/stretchr/testify/assert"
)

// Condition uses a Comparison to assert a complex condition.
func Condition(t TestingT, comp assert.Comparison, msgAndArgs ...interface{}) {
	if !assert.Condition(t, comp, msgAndArgs...) {
		t.FailNow()
	}
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Contains(t, "Hello World", "World")
//    assert.Contains(t, ["Hello", "World"], "World")
//    assert.Contains(t, {"Hello": "World"}, "Hello")
func Contains(t TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if !assert.Contains(t, s, contains, msgAndArgs...) {
		t.FailNow()
	}
}

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func DirExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if !assert.DirExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func ElementsMatch(t TestingT, listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	if !assert.ElementsMatch(t, listA, listB, msgAndArgs...) {
		t.FailNow()
	}
}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Empty(t, obj)
func Empty(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if !assert.Empty(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equal(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.Equal(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualError(t, err,  expectedErrorString)
func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface{}) {
	if !assert.EqualError(t, theError, errString, msgAndArgs...) {
		t.FailNow()
	}
}

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValues(t, uint32(123), int32(123))
func EqualValues(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.EqualValues(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err) {
// 	   assert.Equal(t, expectedError, err)
//   }
func Error(t TestingT, err error, msgAndArgs ...interface{}) {
	if !assert.Error(t, err, msgAndArgs...) {
		t.FailNow()
	}
}

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    assert.Eventually(t, func() bool { return true; }, time.Second, 10*time.Millisecond)
func Eventually(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	if !assert.Eventually(t, condition, waitFor, tick, msgAndArgs...) {
		t.FailNow()
	}
}

// Exactly asserts that two objects are equal in value and type.
//
//    assert.Exactly(t, int32(123), int64(123))
func Exactly(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.Exactly(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// Fail reports a failure through
func Fail(t TestingT, failureMessage string, msgAndArgs ...interface{}) {
	if !assert.Fail(t, failureMessage, msgAndArgs...) {
		t.FailNow()
	}
}

// FailNow fails test
func FailNow(t TestingT, failureMessage string, msgAndArgs ...interface{}) {
	if !assert.FailNow(t, failureMessage, msgAndArgs...) {
		t.FailNow()
	}
}

// False asserts that the specified value is false.
//
//    assert.False(t, myBool)
func False(t TestingT, value bool, msgAndArgs ...interface{}) {
	if !assert.False(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func FileExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if !assert.FileExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}

// Greater asserts that the first element is greater than the second
//
//    assert.Greater(t, 2, 1)
//    assert.Greater(t, float64(2), float64(1))
//    assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if !assert.Greater(t, e1, e2, msgAndArgs...) {
		t.FailNow()
	}
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqual(t, 2, 1)
//    assert.GreaterOrEqual(t, 2, 2)
//    assert.GreaterOrEqual(t, "b", "a")
//    assert.GreaterOrEqual(t, "b", "b")
func GreaterOrEqual(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if !assert.GreaterOrEqual(t, e1, e2, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  assert.HTTPBodyContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if !assert.HTTPBodyContains(t, handler, method, url, values, str, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  assert.HTTPBodyNotContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyNotContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if !assert.HTTPBodyNotContains(t, handler, method, url, values, str, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPError asserts that a specified handler returns an error status code.
//
//  assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPError(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if !assert.HTTPError(t, handler, method, url, values, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPRedirect(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if !assert.HTTPRedirect(t, handler, method, url, values, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
//  assert.HTTPStatusCode(t, myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPStatusCode(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface{}) {
	if !assert.HTTPStatusCode(t, handler, method, url, values, statuscode, msgAndArgs...) {
		t.FailNow()
	}
}

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPSuccess(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if !assert.HTTPSuccess(t, handler, method, url, values, msgAndArgs...) {
		t.FailNow()
	}
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func Implements(t TestingT, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	if !assert.Implements(t, interfaceObject, object, msgAndArgs...) {
		t.FailNow()
	}
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 assert.InDelta(t, math.Pi, 22/7.0, 0.01)
func InDelta(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if !assert.InDelta(t, expected, actual, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func InDeltaMapValues(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if !assert.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func InDeltaSlice(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if !assert.InDeltaSlice(t, expected, actual, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func InEpsilon(t TestingT, expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if !assert.InEpsilon(t, expected, actual, epsilon, msgAndArgs...) {
		t.FailNow()
	}
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlice(t TestingT, expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if !assert.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...) {
		t.FailNow()
	}
}

// IsType asserts that the specified objects are of the same type.
func IsType(t TestingT, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	if !assert.IsType(t, expectedType, object, msgAndArgs...) {
		t.FailNow()
	}
}

// JSONEq asserts that two JSON strings are equivalent.
//
//  assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func JSONEq(t TestingT, expected string, actual string, msgAndArgs ...interface{}) {
	if !assert.JSONEq(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    assert.Len(t, mySlice, 3)
func Len(t TestingT, object interface{}, length int, msgAndArgs ...interface{}) {
	if !assert.Len(t, object, length, msgAndArgs...) {
		t.FailNow()
	}
}

// Less asserts that the first element is less than the second
//
//    assert.Less(t, 1, 2)
//    assert.Less(t, float64(1), float64(2))
//    assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if !assert.Less(t, e1, e2, msgAndArgs...) {
		t.FailNow()
	}
}

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqual(t, 1, 2)
//    assert.LessOrEqual(t, 2, 2)
//    assert.LessOrEqual(t, "a", "b")
//    assert.LessOrEqual(t, "b", "b")
func LessOrEqual(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if !assert.LessOrEqual(t, e1, e2, msgAndArgs...) {
		t.FailNow()
	}
}

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    assert.Never(t, func() bool { return false; }, time.Second, 10*time.Millisecond)
func Never(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	if !assert.Never(t, condition, waitFor, tick, msgAndArgs...) {
		t.FailNow()
	}
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err)
func Nil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if !assert.Nil(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func NoDirExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if !assert.NoDirExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
// 	   assert.Equal(t, expectedObj, actualObj)
//   }
func NoError(t TestingT, err error, msgAndArgs ...interface{}) {
	if !assert.NoError(t, err, msgAndArgs...) {
		t.FailNow()
	}
}

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func NoFileExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if !assert.NoFileExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth")
//    assert.NotContains(t, ["Hello", "World"], "Earth")
//    assert.NotContains(t, {"Hello": "World"}, "Earth")
func NotContains(t TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if !assert.NotContains(t, s, contains, msgAndArgs...) {
		t.FailNow()
	}
}

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmpty(t, obj) {
//    assert.Equal(t, "two", obj[1])
//  }
func NotEmpty(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if !assert.NotEmpty(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.NotEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
//    assert.NotEqualValues(t, obj1, obj2)
func NotEqualValues(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.NotEqualValues(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err)
func NotNil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if !assert.NotNil(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func(){ RemainCalm() })
func NotPanics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !assert.NotPanics(t, f, msgAndArgs...) {
		t.FailNow()
	}
}

// NotRegexp asserts that a specified regexp does not match a string.
//
//  assert.NotRegexp(t, regexp.MustCompile("starts"), "it's starting")
//  assert.NotRegexp(t, "^start", "it's not starting")
func NotRegexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if !assert.NotRegexp(t, rx, str, msgAndArgs...) {
		t.FailNow()
	}
}

// NotSame asserts that two pointers do not reference the same object.
//
//    assert.NotSame(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSame(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.NotSame(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    assert.NotSubset(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func NotSubset(t TestingT, list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if !assert.NotSubset(t, list, subset, msgAndArgs...) {
		t.FailNow()
	}
}

// NotZero asserts that i is not the zero value for its type.
func NotZero(t TestingT, i interface{}, msgAndArgs ...interface{}) {
	if !assert.NotZero(t, i, msgAndArgs...) {
		t.FailNow()
	}
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){ GoCrazy() })
func Panics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !assert.Panics(t, f, msgAndArgs...) {
		t.FailNow()
	}
}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   assert.PanicsWithError(t, "crazy error", func(){ GoCrazy() })
func PanicsWithError(t TestingT, errString string, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !assert.PanicsWithError(t, errString, f, msgAndArgs...) {
		t.FailNow()
	}
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   assert.PanicsWithValue(t, "crazy error", func(){ GoCrazy() })
func PanicsWithValue(t TestingT, expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !assert.PanicsWithValue(t, expected, f, msgAndArgs...) {
		t.FailNow()
	}
}

// Regexp asserts that a specified regexp matches a string.
//
//  assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
//  assert.Regexp(t, "start...$", "it's not starting")
func Regexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if !assert.Regexp(t, rx, str, msgAndArgs...) {
		t.FailNow()
	}
}

// Same asserts that two pointers reference the same object.
//
//    assert.Same(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Same(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.Same(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    assert.Subset(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func Subset(t TestingT, list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if !assert.Subset(t, list, subset, msgAndArgs...) {
		t.FailNow()
	}
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool)
func True(t TestingT, value bool, msgAndArgs ...interface{}) {
	if !assert.True(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)
func WithinDuration(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	if !assert.WithinDuration(t, expected, actual, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// YAMLEq asserts that two YAML strings are equivalent.
func YAMLEq(t TestingT, expected string, actual string, msgAndArgs ...interface{}) {
	if !assert.YAMLEq(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// Zero asserts that i is the zero value for its type.
func Zero(t TestingT, i interface{}, msgAndArgs ...interface{}) {
	if !assert.Zero(t, i, msgAndArgs...) {
		t.FailNow()
	}
}
