package akamai

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func unsetEnvs(t *testing.T) map[string]string {
	configVars := map[string]struct{}{
		"AKAMAI_ACCESS_TOKEN":  {},
		"AKAMAI_CLIENT_TOKEN":  {},
		"AKAMAI_CLIENT_SECRET": {},
		"AKAMAI_HOST":          {},
		"AKAMAI_MAX_BODY":      {},
	}
	existingEnvs := make(map[string]string)

	globalEnvs := os.Environ()
	for _, env := range globalEnvs {
		envKeyValue := strings.SplitN(env, "=", 2)
		if _, ok := configVars[envKeyValue[0]]; ok {
			existingEnvs[envKeyValue[0]] = envKeyValue[1]
		}
	}
	for key := range existingEnvs {
		err := os.Unsetenv(key)
		assert.NoError(t, err)
	}
	return existingEnvs
}

func restoreEnvs(t *testing.T, envs map[string]string) {
	for k, v := range envs {
		err := os.Setenv(k, v)
		assert.NoError(t, err)
	}
}

func TestSetEdgegridEnvs(t *testing.T) {
	tests := map[string]struct {
		givenMap     map[string]interface{}
		givenSection string
		setEnvs      map[string]string
		expectedEnvs map[string]string
	}{
		"no section provided": {
			givenMap: map[string]interface{}{
				"access_token":  "test_access_token",
				"client_token":  "test_client_token",
				"client_secret": "test_client_secret",
				"host":          "test_host",
				"max_body":      123,
			},
			expectedEnvs: map[string]string{
				"AKAMAI_ACCESS_TOKEN":  "test_access_token",
				"AKAMAI_CLIENT_TOKEN":  "test_client_token",
				"AKAMAI_CLIENT_SECRET": "test_client_secret",
				"AKAMAI_HOST":          "test_host",
				"AKAMAI_MAX_BODY":      "123",
			},
		},
		"custom section provided": {
			givenMap: map[string]interface{}{
				"access_token":  "test_access_token",
				"client_token":  "test_client_token",
				"client_secret": "test_client_secret",
				"host":          "test_host",
				"max_body":      123,
			},
			givenSection: "test",
			expectedEnvs: map[string]string{
				"AKAMAI_TEST_ACCESS_TOKEN":  "test_access_token",
				"AKAMAI_TEST_CLIENT_TOKEN":  "test_client_token",
				"AKAMAI_TEST_CLIENT_SECRET": "test_client_secret",
				"AKAMAI_TEST_HOST":          "test_host",
				"AKAMAI_TEST_MAX_BODY":      "123",
			},
		},
		"envs are already set": {
			givenMap: map[string]interface{}{
				"access_token":  "test_access_token",
				"client_token":  "test_client_token",
				"client_secret": "test_client_secret",
				"host":          "test_host",
				"max_body":      123,
			},
			givenSection: "test",
			setEnvs: map[string]string{
				"AKAMAI_TEST_ACCESS_TOKEN":  "existing_access_token",
				"AKAMAI_TEST_CLIENT_TOKEN":  "existing_client_token",
				"AKAMAI_TEST_CLIENT_SECRET": "existing_client_secret",
				"AKAMAI_TEST_HOST":          "existing_host",
				"AKAMAI_TEST_MAX_BODY":      "321",
			},
			expectedEnvs: map[string]string{
				"AKAMAI_TEST_ACCESS_TOKEN":  "existing_access_token",
				"AKAMAI_TEST_CLIENT_TOKEN":  "existing_client_token",
				"AKAMAI_TEST_CLIENT_SECRET": "existing_client_secret",
				"AKAMAI_TEST_HOST":          "existing_host",
				"AKAMAI_TEST_MAX_BODY":      "321",
			},
		},
	}

	existingEnvs := unsetEnvs(t)
	defer restoreEnvs(t, existingEnvs)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			currentEnvs := make(map[string]string, len(test.expectedEnvs))
			for k := range test.expectedEnvs {
				currentEnvs[k] = os.Getenv(k)
				err := os.Unsetenv(k)
				require.NoError(t, err)
			}
			defer func() {
				for k, v := range currentEnvs {
					err := os.Setenv(k, v)
					require.NoError(t, err)
				}
			}()
			for k, v := range test.setEnvs {
				require.NoError(t, os.Setenv(k, v))
			}

			err := setEdgegridEnvs(test.givenMap, test.givenSection)
			require.NoError(t, err)
			for k, v := range test.expectedEnvs {
				assert.Equal(t, v, os.Getenv(k))
			}
		})
	}
}

func TestSetWrongTypeForEdgegridEnvs(t *testing.T) {

	tests := map[string]struct {
		environmentVars map[string]interface{}
	}{
		"nil value for access_token": {
			environmentVars: map[string]interface{}{"access_token": nil},
		},
		"nil value for client_token": {
			environmentVars: map[string]interface{}{"client_token": nil},
		},
		"nil value for host": {
			environmentVars: map[string]interface{}{"host": nil},
		},
		"nil value for client_secret": {
			environmentVars: map[string]interface{}{"client_secret": nil},
		},
		"wrong type of max_body value": {
			environmentVars: map[string]interface{}{"max_body": "not a number"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := setEdgegridEnvs(test.environmentVars, "some section")
			assert.True(t, errors.Is(err, tools.ErrInvalidType))
		})
	}
}

func TestConfigureCache_EnabledInContext(t *testing.T) {
	tests := map[string]struct {
		resourceLocalData   *schema.ResourceData
		expectedMeta        meta
		expectedDiagnostics diag.Diagnostics
	}{
		"cache is enabled": {
			resourceLocalData: getResourceLocalDataWithBoolValue(t, "cache_enabled", true),
			expectedMeta: meta{
				cacheEnabled: true,
			},
			expectedDiagnostics: diag.Diagnostics(nil),
		},
		"cache is not enabled": {
			resourceLocalData: getResourceLocalDataWithBoolValue(t, "cache_enabled", false),
			expectedMeta: meta{
				cacheEnabled: false,
			},
			expectedDiagnostics: diag.Diagnostics(nil),
		},
	}
	for name, test := range tests {
		ctx := context.Background()
		t.Run(name, func(t *testing.T) {

			configuredContext, diagnostics := configureContext(ctx, test.resourceLocalData)
			metaCtx := configuredContext.(*meta)

			assert.Equal(t, test.expectedDiagnostics, diagnostics)
			assert.Equal(t, test.expectedMeta.cacheEnabled, metaCtx.cacheEnabled)
			assert.NotEmpty(t, metaCtx.log)
			assert.NotEmpty(t, metaCtx.operationID)
			assert.NotEmpty(t, metaCtx.sess)
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
			expectedDiagnostics: diag.Errorf(ConfigurationIsNotSpecified),
			withError:           true,
		},
		"with empty edgerc path, default path is used": {
			resourceLocalData:   getResourceLocalData(t, "edgerc", ""),
			expectedDiagnostics: diag.Diagnostics(nil),
			withError:           false,
		},
	}

	existingEnvs := unsetEnvs(t)
	defer restoreEnvs(t, existingEnvs)

	for name, test := range tests {
		ctx := context.Background()
		t.Run(name, func(t *testing.T) {
			meta, diagnostics := configureContext(ctx, test.resourceLocalData)
			if test.withError {
				assert.Nil(t, meta)
			} else {
				assert.NotEmpty(t, meta)
			}
			assert.Equal(t, diagnostics, test.expectedDiagnostics)
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

func TestGetEdgercPath(t *testing.T) {
	tests := map[string]struct {
		edgercPath         string
		expectedEdgercPath string
	}{
		"empty edgerc path": {
			edgercPath:         "",
			expectedEdgercPath: edgegrid.DefaultConfigFile,
		},
		"existing edgerc path": {
			edgercPath:         "../.edgerc",
			expectedEdgercPath: "../.edgerc",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			edgercPath := getEdgercPath(test.edgercPath)
			assert.Equal(t, test.expectedEdgercPath, edgercPath)
		})
	}
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

	resourceDataMap := map[string]interface{}{
		"cache_enabled":  true,
		"edgerc":         "testdata/edgerc",
		"config_section": "validate_edgerc",
	}

	wrongHostFromFile := "\"akaa-ay3i6htctb4uuahh-tklu4vvwja5wzytu.luna-dev.akamaiapis.net/\""

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	configuredContext, diagnostics := configureContext(ctx, resourceData)

	assert.Nil(t, configuredContext)
	assert.Contains(t, diagnostics[0].Summary, wrongHostFromFile)
}

func Test_mergeSchema(t *testing.T) {
	tests := map[string]struct {
		from              map[string]*schema.Schema
		to                map[string]*schema.Schema
		expectedSchemaMap map[string]*schema.Schema
		errorExpected     bool
	}{
		"both schemas are the same": {
			from:          getFirstSchema(),
			to:            getFirstSchema(),
			errorExpected: true,
		},
		"two different schemas": {
			from: getFirstSchema(),
			to:   getSecondSchema(),
			expectedSchemaMap: map[string]*schema.Schema{
				"test element": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "element with type string",
				},
				"different test element": {
					Optional:    true,
					Type:        schema.TypeBool,
					Description: "element with type bool",
				},
			},
			errorExpected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mergedSchema, err := mergeSchema(test.from, test.to)
			if test.errorExpected {
				assert.Nil(t, mergedSchema)
				assert.True(t, errors.Is(err, ErrDuplicateSchemaKey))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedSchemaMap, mergedSchema)
			}
		})
	}
}

func Test_mergeResource(t *testing.T) {
	tests := map[string]struct {
		from                map[string]*schema.Resource
		to                  map[string]*schema.Resource
		expectedResourceMap map[string]*schema.Resource
		errorExpected       bool
	}{
		"both resources are the same": {
			from: map[string]*schema.Resource{
				"some resource": {
					Schema: getFirstSchema(),
				},
			},
			to: map[string]*schema.Resource{
				"some resource": {
					Schema: getFirstSchema(),
				},
			},
			errorExpected: true,
		},
		"two different resources": {
			from: map[string]*schema.Resource{
				"first resource": {
					Schema: getFirstSchema(),
				},
			},
			to: map[string]*schema.Resource{
				"second resource": {
					Schema: getSecondSchema(),
				},
			},
			expectedResourceMap: map[string]*schema.Resource{
				"first resource": {
					Schema: getFirstSchema(),
				},
				"second resource": {
					Schema: getSecondSchema(),
				},
			},
			errorExpected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mergedResource, err := mergeResource(test.from, test.to)
			if test.errorExpected {
				assert.Nil(t, mergedResource)
				assert.True(t, errors.Is(err, ErrDuplicateSchemaKey))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResourceMap, mergedResource)
			}
		})
	}
}

func getFirstSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"test element": {
			Optional:    true,
			Type:        schema.TypeString,
			Description: "element with type string",
		},
	}
}

func getSecondSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"different test element": {
			Optional:    true,
			Type:        schema.TypeBool,
			Description: "element with type bool",
		},
	}
}
