package akamai_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/providers/registry"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestFrameworkProvider(t *testing.T) {
	t.Parallel()
	resp := provider.SchemaResponse{}

	prov := akamai.NewFrameworkProvider(registry.FrameworkSubproviders()...)()
	prov.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	assert.False(t, resp.Diagnostics.HasError())
}

func TestFramework_ConfigureCache_EnabledInContext(t *testing.T) {
	tests := map[string]struct {
		cacheEnabled              bool
		expectedCacheEnabledState bool
	}{
		"cache is enabled": {
			cacheEnabled:              true,
			expectedCacheEnabledState: true,
		},
		"cache is not enabled": {
			cacheEnabled:              false,
			expectedCacheEnabledState: false,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			cache.Enable(false)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV5ProviderFactories: newProtoV5ProviderFactory(dummy{}),
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							provider "akamai" {
								cache_enabled = %v
							}
							data "akamai_dummy" "test" {}
						`, testcase.cacheEnabled),
					},
				},
			})

			assert.Equal(t, testcase.expectedCacheEnabledState, cache.IsEnabled())
		})
	}
}

func TestFramework_ConfigureEdgercInContext(t *testing.T) {
	tests := map[string]struct {
		key           string
		value         string
		expectedError *regexp.Regexp
	}{
		"file with EdgeGrid configuration does not exist": {
			key:           "edgerc",
			value:         "not_existing_file_path",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: loading config file: open\nnot_existing_file_path: no such file or directory"),
		},
		"config section does not exist": {
			key:           "config_section",
			value:         "not_existing_config_section",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: provided config section does not\nexist: section \"not_existing_config_section\" does not exist"),
		},
		"with empty edgerc path, default path is used": {
			key:           "edgerc",
			value:         "",
			expectedError: nil,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV5ProviderFactories: newProtoV5ProviderFactory(dummy{}),
				Steps: []resource.TestStep{
					{
						ExpectError: testcase.expectedError,
						Config: fmt.Sprintf(`
							provider "akamai" {
								%v = "%v"
							}
							data "akamai_dummy" "test" {}
						`, testcase.key, testcase.value),
					},
				},
			})
		})
	}
}

func TestFramework_EdgercValidate(t *testing.T) {

	tests := map[string]struct {
		expectedError *regexp.Regexp
		configSection string
	}{
		"no host": {
			configSection: "no_host",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from\nedgerc: \"host\""),
		},
		"no client_secret": {
			configSection: "no_client_secret",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from\nedgerc: \"client_secret\""),
		},
		"no access_token": {
			configSection: "no_access_token",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from\nedgerc: \"access_token\""),
		},
		"no client_token": {
			configSection: "no_client_token",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from\nedgerc: \"client_token\""),
		},
		"wrong format of host": {
			configSection: "validate_edgerc",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: host must not contain '/' at the\nend: \"host.com/\""),
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV5ProviderFactories: newProtoV5ProviderFactory(dummy{}),
				Steps: []resource.TestStep{
					{
						ExpectError: testcase.expectedError,
						Config: fmt.Sprintf(`
							provider "akamai" {
								cache_enabled  = true
								edgerc         = "testdata/edgerc"
								config_section = "%v"
							}
							data "akamai_dummy" "test" {}
					`, testcase.configSection),
					},
				},
			})
		})
	}
}

func newProtoV5ProviderFactory(subproviders ...subprovider.Framework) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"akamai": func() (tfprotov5.ProviderServer, error) {
			return providerserver.NewProtocol5(
				akamai.NewFrameworkProvider(subproviders...)(),
			)(), nil
		},
	}
}
