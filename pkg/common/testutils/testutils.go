// Package testutils gathers reusable pieces useful for testing
package testutils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// tfTestTempDir specifies the location of tmp directory which will be used by provider SDK's testing framework
const tfTestTempDir = "./test_tmp"

// TestRunner executes common test setup and teardown in all subproviders
func TestRunner(m *testing.M) {
	if err := TFTestSetup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := TFTestTeardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

// TFTestSetup contains common setup for tests in all subproviders
func TFTestSetup() error {
	if err := os.MkdirAll(tfTestTempDir, 0755); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	if err := os.Setenv("TF_ACC_TEMP_DIR", tfTestTempDir); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	return nil
}

// TFTestTeardown contains common teardown for tests in all subproviders
func TFTestTeardown() error {
	if err := os.RemoveAll(tfTestTempDir); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	if err := os.Unsetenv("TF_ACC_TEMP_DIR"); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	return nil
}

// TattleT wraps a *testing.T to intercept a Testify mock's call of t.FailNow(). When testing.t.FailNow() is called from
// any goroutine other than the one on which a test was created, it causes the test to hang. Testify's mocks fail to
// inform the user which test failed. Use this struct to wrap a *testing.TattleT when you call `mock.Test(TattleT{t})`
// and the mock's failure message will include the failling test's name. Such failures are usually caused by unexpected
// method calls on a mock.
//
// NB: You would only need to use this where Testify mocks are used in tests that spawn goroutines, such as those run by
// the Terraform test driver.
type TattleT struct{ *testing.T }

// FailNow overrides testing.T.FailNow() so when a test mock fails an assertion, we see which test failed
func (t TattleT) FailNow() {
	t.T.Helper()
	t.T.Fatalf("FAIL: %s", t.T.Name())
}

// TestStepDestroyFailed creates a terraform test step, to work around an issue of verifying error if resource destroy failed correctly.
// Problem: When throwing an error from destroy terraform units produce “Error running post-test destroy, there may be dangling resources" because
// Terraform requires that destroy is successful to clean up the resources. Adding this step to the end will cause destroy to be triggered twice:
// First one will test the erroneous behaviour, the second one will be successful (to fulfil Terraform units requirement).
// Usage: add returned test step as the last step with the same config file as in previous step and expected error message.
func TestStepDestroyFailed(config string, expectedError *regexp.Regexp) resource.TestStep {
	return resource.TestStep{
		Config:      config,
		Destroy:     true,
		ExpectError: expectedError,
	}
}
