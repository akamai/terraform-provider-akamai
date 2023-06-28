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
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
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
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: loading config file: open not_existing_file_path: no such file or directory"),
		},
		"config section does not exist": {
			key:           "config_section",
			value:         "not_existing_config_section",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: provided config section does not exist: section \"not_existing_config_section\" does not exist"),
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
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from edgerc: \"host\""),
		},
		"no client_secret": {
			configSection: "no_client_secret",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from edgerc: \"client_secret\""),
		},
		"no access_token": {
			configSection: "no_access_token",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from edgerc: \"access_token\""),
		},
		"no client_token": {
			configSection: "no_client_token",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: required option is missing from edgerc: \"client_token\""),
		},
		"wrong format of host": {
			configSection: "validate_edgerc",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: host must not contain '/' at the end: \"host.com/\""),
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

func TestFramework_EdgercFromConfig(t *testing.T) {
	tests := map[string]struct {
		expectedError *regexp.Regexp
		clientSecret  string
		host          string
		accessToken   string
		clientToken   string
	}{
		"valid config": {
			host:         "host.com",
			clientSecret: "client_secret",
			accessToken:  "access_token",
			clientToken:  "client_token",
		},
		"invalid - empty host": {
			clientSecret:  "client_secret",
			accessToken:   "access_token",
			clientToken:   "client_token",
			expectedError: regexp.MustCompile("Attribute host cannot be empty"),
		},
		"invalid - empty client_secret": {
			host:          "host.com",
			accessToken:   "access_token",
			clientToken:   "client_token",
			expectedError: regexp.MustCompile("Attribute client_secret cannot be empty"),
		},
		"invalid - empty access_token": {
			host:          "host.com",
			clientSecret:  "client_secret",
			clientToken:   "client_token",
			expectedError: regexp.MustCompile("Attribute access_token cannot be empty"),
		},
		"invalid - empty client_token": {
			clientSecret:  "client_secret",
			host:          "host.com",
			accessToken:   "access_token",
			expectedError: regexp.MustCompile("Attribute client_token cannot be empty"),
		},
		"wrong format of host": {
			clientSecret:  "client_secret",
			host:          "host.com/",
			accessToken:   "access_token",
			clientToken:   "client_token",
			expectedError: regexp.MustCompile("error reading Akamai EdgeGrid configuration: host must not contain '/' at the end: \"host.com/\""),
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
								config {
									client_secret = "%v"
    								host          = "%v"
    								access_token  = "%v"
    								client_token  = "%v"
								}
							}
							data "akamai_dummy" "test" {}
					`, testcase.clientSecret, testcase.host, testcase.accessToken, testcase.clientToken),
					},
				},
			})
		})
	}
}

func TestFramework_EdgercFromConfig_missing_required_attributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV5ProviderFactories: newProtoV5ProviderFactory(dummy{}),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("The argument \"host\" is required, but no definition was found"),
				Config: `
					provider "akamai" {
						config {
							client_secret = "client_secret"
							access_token  = "access_token"
							client_token  = "client_token"
						}
					}
					data "akamai_dummy" "test" {}`,
			},
			{
				ExpectError: regexp.MustCompile("The argument \"client_secret\" is required, but no definition was found"),
				Config: `
					provider "akamai" {
						config {
							host          = "host"
							access_token  = "access_token"
							client_token  = "client_token"
						}
					}
					data "akamai_dummy" "test" {}`,
			},
			{
				ExpectError: regexp.MustCompile("The argument \"access_token\" is required, but no definition was found"),
				Config: `
					provider "akamai" {
						config {
							host          = "host"
							client_secret = "client_secret"
							client_token  = "client_token"
						}
					}
					data "akamai_dummy" "test" {}`,
			},
			{
				ExpectError: regexp.MustCompile("The argument \"client_token\" is required, but no definition was found"),
				Config: `
					provider "akamai" {
						config {
							host          = "host"
							client_secret = "client_secret"
							access_token  = "access_token"
						}
					}
					data "akamai_dummy" "test" {}`,
			},
		},
	})
}

func newProtoV5ProviderFactory(subproviders ...subprovider.Framework) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"akamai": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				akamai.NewPluginProvider()().GRPCProvider,
				providerserver.NewProtocol5(
					akamai.NewFrameworkProvider(subproviders...)(),
				),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
