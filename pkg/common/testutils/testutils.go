// Package testutils gathers reusable pieces useful for testing
package testutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/stretchr/testify/require"
)

// tfTestTempDir specifies the location of tmp directory which will be used by provider SDK's testing framework
const tfTestTempDir = "./test_tmp"

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

// NewSDKProviderFactories uses provided SDK subprovider to create provider factories for test purposes
func NewSDKProviderFactories(subprovider subprovider.SDK) map[string]func() (tfprotov6.ProviderServer, error) {
	testAccSDKProvider := akamai.NewSDKProvider(subprovider)()
	testAccProviders := map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			sdkProviderV6, err := tf5to6server.UpgradeServer(
				context.Background(),
				testAccSDKProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return sdkProviderV6
				},
			}

			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
	return testAccProviders
}
