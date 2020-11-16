package akamai

import (
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"os"
	"testing"
)

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
