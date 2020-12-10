package test

import (
	"testing"
)

// TattleT wraps a *testing.TattleT to intercept a Testify mock's call of t.FailNow(). When testing.t.FailNow() is
// called from any goroutine other than the one on which a test was created, it causes the test to hang. Testify's mocks
// fail to inform the user which test failed. Use this struct to wrap a *testing.TattleT when you call
// `mock.Test(TattleT{t})` and the mock's failure message will include the failling test's name. Such failures are
// usually caused by unexpected methodcalls on a mock.
//
// NB: You would only need to use this where Testify mocks are used in tests that spawn goroutines, such as those run by
//     the Terraform test driver.
type TattleT struct{ *testing.T }

// FailNow overrides testing.T.FailNow() so when a test mock fails an assertion, we see which test failed
func (t TattleT) FailNow() {
	t.T.Fatalf("FAIL: %s", t.T.Name())
}
