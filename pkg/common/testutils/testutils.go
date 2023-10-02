// Package testutils gathers reusable pieces useful for testing
package testutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/stretchr/testify/require"
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

// LoadFixtureBytes returns the entire contents of the given file as a byte slice
func LoadFixtureBytes(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return contents
}

// LoadFixtureString returns the entire contents of the given file as a string
func LoadFixtureString(t *testing.T, format string, args ...interface{}) string {
	return string(LoadFixtureBytes(t, fmt.Sprintf(format, args...)))
}

// NewProtoV6ProviderFactory uses provided subprovider to create provider factory for test purposes
func NewProtoV6ProviderFactory(subproviders ...subprovider.Subprovider) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			sdkProviderV6, err := akamai.NewProtoV6SDKProvider(subproviders)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				sdkProviderV6,
				providerserver.NewProtocol6(
					akamai.NewFrameworkProvider(subproviders...)(),
				),
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
