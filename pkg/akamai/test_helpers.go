package akamai

import (
	"fmt"
	"os"
)

const TFTestTempDir = "./test_tmp"

func TFTestSetup() error {
	if err := os.MkdirAll(TFTestTempDir, 0755); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	if err := os.Setenv("TF_ACC_TEMP_DIR", TFTestTempDir); err != nil {
		return fmt.Errorf("test setup failed: %s", err)
	}
	return nil
}

func TFTestTeardown() error {
	if err := os.RemoveAll(TFTestTempDir); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	if err := os.Unsetenv("TF_ACC_TEMP_DIR"); err != nil {
		return fmt.Errorf("test teardown failed: %s", err)
	}
	return nil
}
