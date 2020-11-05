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
func (a *Assertions) Condition(comp assert.Comparison, msgAndArgs ...interface{}) {
	Condition(a.t, comp, msgAndArgs...)
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    a.Contains("Hello World", "World")
//    a.Contains(["Hello", "World"], "World")
//    a.Contains({"Hello": "World"}, "Hello")
func (a *Assertions) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	Contains(a.t, s, contains, msgAndArgs...)
}

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func (a *Assertions) DirExists(path string, msgAndArgs ...interface{}) {
	DirExists(a.t, path, msgAndArgs...)
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// a.ElementsMatch([1, 3, 2, 3], [1, 3, 3, 2])
func (a *Assertions) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	ElementsMatch(a.t, listA, listB, msgAndArgs...)
}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  a.Empty(obj)
func (a *Assertions) Empty(object interface{}, msgAndArgs ...interface{}) {
	Empty(a.t, object, msgAndArgs...)
}

// Equal asserts that two objects are equal.
//
//    a.Equal(123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func (a *Assertions) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	Equal(a.t, expected, actual, msgAndArgs...)
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   a.EqualError(err,  expectedErrorString)
func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface{}) {
	EqualError(a.t, theError, errString, msgAndArgs...)
}

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    a.EqualValues(uint32(123), int32(123))
func (a *Assertions) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	EqualValues(a.t, expected, actual, msgAndArgs...)
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.Error(err) {
// 	   assert.Equal(t, expectedError, err)
//   }
func (a *Assertions) Error(err error, msgAndArgs ...interface{}) {
	Error(a.t, err, msgAndArgs...)
}

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    a.Eventually(func() bool { return true; }, time.Second, 10*time.Millisecond)
func (a *Assertions) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	Eventually(a.t, condition, waitFor, tick, msgAndArgs...)
}

// Exactly asserts that two objects are equal in value and type.
//
//    a.Exactly(int32(123), int64(123))
func (a *Assertions) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	Exactly(a.t, expected, actual, msgAndArgs...)
}

// Fail reports a failure through
func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface{}) {
	Fail(a.t, failureMessage, msgAndArgs...)
}

// FailNow fails test
func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface{}) {
	FailNow(a.t, failureMessage, msgAndArgs...)
}

// False asserts that the specified value is false.
//
//    a.False(myBool)
func (a *Assertions) False(value bool, msgAndArgs ...interface{}) {
	False(a.t, value, msgAndArgs...)
}

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (a *Assertions) FileExists(path string, msgAndArgs ...interface{}) {
	FileExists(a.t, path, msgAndArgs...)
}

// Greater asserts that the first element is greater than the second
//
//    a.Greater(2, 1)
//    a.Greater(float64(2), float64(1))
//    a.Greater("b", "a")
func (a *Assertions) Greater(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	Greater(a.t, e1, e2, msgAndArgs...)
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    a.GreaterOrEqual(2, 1)
//    a.GreaterOrEqual(2, 2)
//    a.GreaterOrEqual("b", "a")
//    a.GreaterOrEqual("b", "b")
func (a *Assertions) GreaterOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	GreaterOrEqual(a.t, e1, e2, msgAndArgs...)
}

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  a.HTTPBodyContains(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
}

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  a.HTTPBodyNotContains(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
}

// HTTPError asserts that a specified handler returns an error status code.
//
//  a.HTTPError(myHandler, "POST", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPError(a.t, handler, method, url, values, msgAndArgs...)
}

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  a.HTTPRedirect(myHandler, "GET", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
}

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
//  a.HTTPStatusCode(myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPStatusCode(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface{}) {
	HTTPStatusCode(a.t, handler, method, url, values, statuscode, msgAndArgs...)
}

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  a.HTTPSuccess(myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
}

// Implements asserts that an object is implemented by the specified interface.
//
//    a.Implements((*MyInterface)(nil), new(MyObject))
func (a *Assertions) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	Implements(a.t, interfaceObject, object, msgAndArgs...)
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 a.InDelta(math.Pi, 22/7.0, 0.01)
func (a *Assertions) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDelta(a.t, expected, actual, delta, msgAndArgs...)
}

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (a *Assertions) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func (a *Assertions) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func (a *Assertions) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func (a *Assertions) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
}

// IsType asserts that the specified objects are of the same type.
func (a *Assertions) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	IsType(a.t, expectedType, object, msgAndArgs...)
}

// JSONEq asserts that two JSON strings are equivalent.
//
//  a.JSONEq(`{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface{}) {
	JSONEq(a.t, expected, actual, msgAndArgs...)
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    a.Len(mySlice, 3)
func (a *Assertions) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	Len(a.t, object, length, msgAndArgs...)
}

// Less asserts that the first element is less than the second
//
//    a.Less(1, 2)
//    a.Less(float64(1), float64(2))
//    a.Less("a", "b")
func (a *Assertions) Less(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	Less(a.t, e1, e2, msgAndArgs...)
}

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    a.LessOrEqual(1, 2)
//    a.LessOrEqual(2, 2)
//    a.LessOrEqual("a", "b")
//    a.LessOrEqual("b", "b")
func (a *Assertions) LessOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	LessOrEqual(a.t, e1, e2, msgAndArgs...)
}

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    a.Never(func() bool { return false; }, time.Second, 10*time.Millisecond)
func (a *Assertions) Never(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	Never(a.t, condition, waitFor, tick, msgAndArgs...)
}

// Nil asserts that the specified object is nil.
//
//    a.Nil(err)
func (a *Assertions) Nil(object interface{}, msgAndArgs ...interface{}) {
	Nil(a.t, object, msgAndArgs...)
}

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func (a *Assertions) NoDirExists(path string, msgAndArgs ...interface{}) {
	NoDirExists(a.t, path, msgAndArgs...)
}

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.NoError(err) {
// 	   assert.Equal(t, expectedObj, actualObj)
//   }
func (a *Assertions) NoError(err error, msgAndArgs ...interface{}) {
	NoError(a.t, err, msgAndArgs...)
}

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func (a *Assertions) NoFileExists(path string, msgAndArgs ...interface{}) {
	NoFileExists(a.t, path, msgAndArgs...)
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    a.NotContains("Hello World", "Earth")
//    a.NotContains(["Hello", "World"], "Earth")
//    a.NotContains({"Hello": "World"}, "Earth")
func (a *Assertions) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	NotContains(a.t, s, contains, msgAndArgs...)
}

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if a.NotEmpty(obj) {
//    assert.Equal(t, "two", obj[1])
//  }
func (a *Assertions) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	NotEmpty(a.t, object, msgAndArgs...)
}

// NotEqual asserts that the specified values are NOT equal.
//
//    a.NotEqual(obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (a *Assertions) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	NotEqual(a.t, expected, actual, msgAndArgs...)
}

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
//    a.NotEqualValues(obj1, obj2)
func (a *Assertions) NotEqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	NotEqualValues(a.t, expected, actual, msgAndArgs...)
}

// NotNil asserts that the specified object is not nil.
//
//    a.NotNil(err)
func (a *Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) {
	NotNil(a.t, object, msgAndArgs...)
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   a.NotPanics(func(){ RemainCalm() })
func (a *Assertions) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	NotPanics(a.t, f, msgAndArgs...)
}

// NotRegexp asserts that a specified regexp does not match a string.
//
//  a.NotRegexp(regexp.MustCompile("starts"), "it's starting")
//  a.NotRegexp("^start", "it's not starting")
func (a *Assertions) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	NotRegexp(a.t, rx, str, msgAndArgs...)
}

// NotSame asserts that two pointers do not reference the same object.
//
//    a.NotSame(ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) NotSame(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	NotSame(a.t, expected, actual, msgAndArgs...)
}

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    a.NotSubset([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func (a *Assertions) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	NotSubset(a.t, list, subset, msgAndArgs...)
}

// NotZero asserts that i is not the zero value for its type.
func (a *Assertions) NotZero(i interface{}, msgAndArgs ...interface{}) {
	NotZero(a.t, i, msgAndArgs...)
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   a.Panics(func(){ GoCrazy() })
func (a *Assertions) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	Panics(a.t, f, msgAndArgs...)
}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   a.PanicsWithError("crazy error", func(){ GoCrazy() })
func (a *Assertions) PanicsWithError(errString string, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	PanicsWithError(a.t, errString, f, msgAndArgs...)
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   a.PanicsWithValue("crazy error", func(){ GoCrazy() })
func (a *Assertions) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	PanicsWithValue(a.t, expected, f, msgAndArgs...)
}

// Regexp asserts that a specified regexp matches a string.
//
//  a.Regexp(regexp.MustCompile("start"), "it's starting")
//  a.Regexp("start...$", "it's not starting")
func (a *Assertions) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	Regexp(a.t, rx, str, msgAndArgs...)
}

// Same asserts that two pointers reference the same object.
//
//    a.Same(ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) Same(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	Same(a.t, expected, actual, msgAndArgs...)
}

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    a.Subset([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func (a *Assertions) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	Subset(a.t, list, subset, msgAndArgs...)
}

// True asserts that the specified value is true.
//
//    a.True(myBool)
func (a *Assertions) True(value bool, msgAndArgs ...interface{}) {
	True(a.t, value, msgAndArgs...)
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   a.WithinDuration(time.Now(), time.Now(), 10*time.Second)
func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
}

// YAMLEq asserts that two YAML strings are equivalent.
func (a *Assertions) YAMLEq(expected string, actual string, msgAndArgs ...interface{}) {
	YAMLEq(a.t, expected, actual, msgAndArgs...)
}

// Zero asserts that i is the zero value for its type.
func (a *Assertions) Zero(i interface{}, msgAndArgs ...interface{}) {
	Zero(a.t, i, msgAndArgs...)
}
