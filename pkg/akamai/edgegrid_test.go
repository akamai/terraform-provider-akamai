package akamai

import (
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEdgegridConfig(t *testing.T) {
	edgercPath := "testdata/edgerc"
	section := "default"
	clientSecret := "test_client_secret"
	clientToken := "test_client_token"
	accessToken := "test_access_token"

	envHost := "env.com"
	configHost := "config.com"
	fileHost := "host.com"

	config := configBearer{
		host:         configHost,
		accessToken:  accessToken,
		clientToken:  clientToken,
		clientSecret: clientSecret,
	}

	t.Run("env is prioritized over config", func(t *testing.T) {
		t.Setenv("AKAMAI_HOST", envHost)
		t.Setenv("AKAMAI_ACCESS_TOKEN", accessToken)
		t.Setenv("AKAMAI_CLIENT_TOKEN", clientToken)
		t.Setenv("AKAMAI_CLIENT_SECRET", clientSecret)

		edgegridConfig, err := newEdgegridConfig("", "", config)
		require.NoError(t, err)
		assert.Equal(t, envHost, edgegridConfig.Host)
	})

	t.Run("env is prioritized over file", func(t *testing.T) {
		t.Setenv("AKAMAI_HOST", envHost)
		t.Setenv("AKAMAI_ACCESS_TOKEN", accessToken)
		t.Setenv("AKAMAI_CLIENT_TOKEN", clientToken)
		t.Setenv("AKAMAI_CLIENT_SECRET", clientSecret)

		edgegridConfig, err := newEdgegridConfig(edgercPath, section, configBearer{})
		require.NoError(t, err)
		assert.Equal(t, envHost, edgegridConfig.Host)
	})

	t.Run("non-default section is used for reading env", func(t *testing.T) {
		testSection := "TEST"
		host := "testenv.com"

		t.Setenv(fmt.Sprintf("AKAMAI_%s_HOST", testSection), host)
		t.Setenv(fmt.Sprintf("AKAMAI_%s_ACCESS_TOKEN", testSection), accessToken)
		t.Setenv(fmt.Sprintf("AKAMAI_%s_CLIENT_TOKEN", testSection), clientToken)
		t.Setenv(fmt.Sprintf("AKAMAI_%s_CLIENT_SECRET", testSection), clientSecret)

		edgegridConfig, err := newEdgegridConfig("", testSection, config)
		require.NoError(t, err)
		assert.Equal(t, host, edgegridConfig.Host)
	})

	t.Run("uses config when provided and env not set", func(t *testing.T) {
		edgegridConfig, err := newEdgegridConfig("", "", config)
		require.NoError(t, err)
		assert.Equal(t, configHost, edgegridConfig.Host)
	})

	t.Run("uses config when provided and env not valid", func(t *testing.T) {
		t.Setenv("AKAMAI_HOST", "env.com")
		t.Setenv("AKAMAI_ACCESS_TOKEN", accessToken)

		edgegridConfig, err := newEdgegridConfig("", "", config)
		require.NoError(t, err)
		assert.Equal(t, configHost, edgegridConfig.Host)
	})

	t.Run("config is prioritized over edgerc file", func(t *testing.T) {
		edgegridConfig, err := newEdgegridConfig(edgercPath, section, config)
		require.NoError(t, err)
		assert.Equal(t, configHost, edgegridConfig.Host)
	})

	t.Run("uses edgerc file when env and config not provided", func(t *testing.T) {
		edgegridConfig, err := newEdgegridConfig(edgercPath, section, configBearer{})
		require.NoError(t, err)
		assert.Equal(t, fileHost, edgegridConfig.Host)
	})

	t.Run("uses edgerc file when provided and env is invalid", func(t *testing.T) {
		t.Setenv("AKAMAI_HOST", "env.com")
		t.Setenv("AKAMAI_ACCESS_TOKEN", accessToken)

		edgegridConfig, err := newEdgegridConfig(edgercPath, section, configBearer{})
		require.NoError(t, err)
		assert.Equal(t, fileHost, edgegridConfig.Host)

	})

	t.Run("uses default edgerc path and section when none provided", func(t *testing.T) {
		edgegridConfig, err := newEdgegridConfig("", "", configBearer{})
		require.NoError(t, err)
		assert.Equal(t, fileHost, edgegridConfig.Host)
	})
}

func TestEdgercPathOrDefault(t *testing.T) {
	t.Parallel()

	path := "testdata/edgerc"

	assert.Equal(t, path, edgercPathOrDefault(path))
	assert.Equal(t, DefaultConfigFilePath, edgercPathOrDefault(""))
}

func TestEdgercSectionOrDefault(t *testing.T) {
	t.Parallel()

	section := "not_default"

	assert.Equal(t, section, edgercSectionOrDefault(section))
	assert.Equal(t, edgegrid.DefaultSection, edgercSectionOrDefault(""))
}

func TestConfigBearerToEdgegridConfig(t *testing.T) {
	t.Parallel()

	accessToken := "test_access_token"
	accountKey := "test_account_key"
	clientSecret := "test_client_secret"
	clientToken := "test_client_token"
	host := "host.com"
	maxBody := 1234

	t.Run("returns error when config not valid", func(t *testing.T) {
		configBearer := configBearer{
			host:        host,
			accessToken: accessToken,
		}

		_, err := configBearer.toEdgegridConfig()
		assert.Error(t, err)
	})
	t.Run("creates edgegrid.Config using all values", func(t *testing.T) {
		configBearer := configBearer{
			accessToken:  accessToken,
			accountKey:   accountKey,
			clientSecret: clientSecret,
			clientToken:  clientToken,
			host:         host,
			maxBody:      maxBody,
		}

		edgegridConfig, err := configBearer.toEdgegridConfig()
		require.NoError(t, err)

		assert.Equal(t, configBearer.accessToken, edgegridConfig.AccessToken)
		assert.Equal(t, configBearer.accountKey, edgegridConfig.AccountKey)
		assert.Equal(t, configBearer.clientSecret, edgegridConfig.ClientSecret)
		assert.Equal(t, configBearer.clientToken, edgegridConfig.ClientToken)
		assert.Equal(t, configBearer.host, edgegridConfig.Host)
		assert.Equal(t, configBearer.maxBody, edgegridConfig.MaxBody)
	})
	t.Run("sets MaxBody to MaxBodySize when 0", func(t *testing.T) {
		configBearer := configBearer{
			accessToken:  accessToken,
			clientSecret: clientSecret,
			clientToken:  clientToken,
			host:         host,
		}

		edgegridConfig, err := configBearer.toEdgegridConfig()
		require.NoError(t, err)

		assert.Equal(t, edgegrid.MaxBodySize, edgegridConfig.MaxBody)
	})
}

func TestConfigBearerValid(t *testing.T) {
	t.Parallel()

	accessToken := "test_access_token"
	clientSecret := "test_client_secret"
	clientToken := "test_client_token"
	host := "host.com"

	testCases := map[string]struct {
		config   configBearer
		expected bool
	}{
		"is valid when only required provided": {
			config: configBearer{
				accessToken:  accessToken,
				clientSecret: clientSecret,
				clientToken:  clientToken,
				host:         host,
			},
			expected: true,
		},
		"is valid when all attributes provided": {
			config: configBearer{
				accessToken:  accessToken,
				clientSecret: clientSecret,
				clientToken:  clientToken,
				host:         host,
				maxBody:      1234,
				accountKey:   "test_account_key",
			},
			expected: true,
		},
		"not valid - empty host": {
			config: configBearer{
				accessToken:  accessToken,
				clientSecret: clientSecret,
				clientToken:  clientToken,
			},
			expected: false,
		},
		"not valid - empty client token": {
			config: configBearer{
				accessToken:  accessToken,
				clientSecret: clientSecret,
				host:         host,
			},
			expected: false,
		},
		"not valid - empty client secret": {
			config: configBearer{
				accessToken: accessToken,
				clientToken: clientToken,
				host:        host,
			},
			expected: false,
		},
		"not valid - empty access token": {
			config: configBearer{
				clientSecret: clientSecret,
				clientToken:  clientToken,
				host:         host,
			},
			expected: false,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.config.valid())
		})
	}
}
