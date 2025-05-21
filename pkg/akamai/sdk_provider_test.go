package akamai_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/cache"

	// Load the providers
	_ "github.com/akamai/terraform-provider-akamai/v8/pkg/providers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/providers/registry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSDKProvider(t *testing.T) {
	t.Parallel()

	provider := akamai.NewSDKProvider(registry.Subproviders()...)()
	err := provider.InternalValidate()
	assert.NoError(t, err)
}

func TestConfigureCache_EnabledInContext(t *testing.T) {
	tests := map[string]struct {
		resourceLocalData         *schema.ResourceData
		expectedCacheEnabledState bool
	}{
		"cache is enabled": {
			resourceLocalData:         getResourceLocalDataWithBoolValue(t, "cache_enabled", true),
			expectedCacheEnabledState: true,
		},
		"cache is not enabled": {
			resourceLocalData:         getResourceLocalDataWithBoolValue(t, "cache_enabled", false),
			expectedCacheEnabledState: false,
		},
	}
	for name, test := range tests {
		ctx := context.Background()
		t.Run(name, func(t *testing.T) {
			prov := akamai.NewSDKProvider()
			_, diagnostics := prov().ConfigureContextFunc(ctx, test.resourceLocalData)
			require.False(t, diagnostics.HasError(), fmt.Sprintf("unexpected error in diagnostics: %v", diagnostics))

			assert.Equal(t, test.expectedCacheEnabledState, cache.IsEnabled())
		})
	}
}

func TestConfigureEdgercInContext(t *testing.T) {
	tests := map[string]struct {
		resourceLocalData   *schema.ResourceData
		expectedDiagnostics diag.Diagnostics
		withError           bool
	}{
		"file with EdgeGrid configuration does not exist": {
			resourceLocalData:   getResourceLocalData(t, "edgerc", "not_existing_file_path"),
			expectedDiagnostics: diag.Errorf("%s: %s: %s", akamai.ErrWrongEdgeGridConfiguration, edgegrid.ErrLoadingFile, "open not_existing_file_path: no such file or directory"),
			withError:           true,
		},
		"config section does not exist": {
			resourceLocalData:   getResourceLocalData(t, "config_section", "not_existing_config_section"),
			expectedDiagnostics: diag.Errorf("%s: %s: %s", akamai.ErrWrongEdgeGridConfiguration, edgegrid.ErrSectionDoesNotExist, "section \"not_existing_config_section\" does not exist"),
			withError:           true,
		},
		"uses defaults with empty edgerc and config_section": {
			resourceLocalData:   getResourceLocalData(t, "edgerc", ""),
			expectedDiagnostics: diag.Diagnostics(nil),
			withError:           false,
		},
	}

	for name, test := range tests {
		ctx := context.Background()
		t.Run(name, func(t *testing.T) {
			prov := akamai.NewSDKProvider()
			meta, diagnostics := prov().ConfigureContextFunc(ctx, test.resourceLocalData)

			if test.withError {
				assert.Nil(t, meta)
			} else {
				assert.NotEmpty(t, meta)
			}
			assert.Equal(t, test.expectedDiagnostics, diagnostics)
		})
	}
}

func getResourceLocalDataWithBoolValue(t *testing.T, key string, value bool) *schema.ResourceData {
	resourceSchema := map[string]*schema.Schema{
		key: {
			Type: schema.TypeBool,
		},
	}
	resourceDataMap := map[string]interface{}{
		key: value,
	}
	return schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
}

func getResourceLocalData(t *testing.T, key string, value interface{}) *schema.ResourceData {
	resourceSchema := map[string]*schema.Schema{
		"cache_enabled": {
			Type: schema.TypeBool,
		},
		key: {
			Type: schema.TypeString,
		},
	}

	dataMap := map[string]interface{}{
		"cache_enabled": true,
		key:             value,
	}
	return schema.TestResourceDataRaw(t, resourceSchema, dataMap)
}

func TestEdgercValidate(t *testing.T) {
	ctx := context.Background()
	resourceSchema := map[string]*schema.Schema{
		"cache_enabled": {
			Type: schema.TypeBool,
		},
		"edgerc": {
			Type: schema.TypeString,
		},
		"config_section": {
			Type: schema.TypeString,
		},
	}
	tests := map[string]struct {
		expectedError error
		configSection string
	}{
		"no host": {
			configSection: "no_host",
			expectedError: fmt.Errorf("%s: \"host\"", edgegrid.ErrRequiredOptionEdgerc),
		},
		"no client_secret": {
			configSection: "no_client_secret",
			expectedError: fmt.Errorf("%s: \"client_secret\"", edgegrid.ErrRequiredOptionEdgerc),
		},
		"no access_token": {
			configSection: "no_access_token",
			expectedError: fmt.Errorf("%s: \"access_token\"", edgegrid.ErrRequiredOptionEdgerc),
		},
		"no client_token": {
			configSection: "no_client_token",
			expectedError: fmt.Errorf("%s: \"client_token\"", edgegrid.ErrRequiredOptionEdgerc),
		},
		"wrong format of host": {
			configSection: "validate_edgerc",
			expectedError: fmt.Errorf("host must not contain '/' at the end: \"host.com/\""),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resourceDataMap := map[string]interface{}{
				"cache_enabled":  true,
				"edgerc":         "testdata/edgerc",
				"config_section": test.configSection,
			}
			resourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

			prov := akamai.NewSDKProvider()
			configuredContext, diagnostics := prov().ConfigureContextFunc(ctx, resourceData)

			assert.Nil(t, configuredContext)
			assert.Contains(t, diagnostics[0].Summary, test.expectedError.Error())
		})
	}
}
