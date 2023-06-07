package akamai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEdgegridConfig(t *testing.T) {
	t.Parallel()

	path := "testdata/edgerc"
	section := "default"
	config := func() map[string]any {
		return map[string]any{
			"host":          "host.com",
			"access_token":  "access_token",
			"client_token":  "client_token",
			"client_secret": "client_secret",
			"max_body":      0,
			"account_key":   "",
		}
	}

	t.Run("from config map", func(t *testing.T) {
		t.Parallel()

		_, err := newEdgegridConfig("", "", config())
		require.NoError(t, err)
	})

	t.Run("from file", func(t *testing.T) {
		t.Parallel()

		_, err := newEdgegridConfig(path, section, nil)
		require.NoError(t, err)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		t.Parallel()

		_, err := newEdgegridConfig(path, "", config())
		assert.Error(t, err)

		_, err = newEdgegridConfig("", section, config())
		assert.Error(t, err)

		_, err = newEdgegridConfig(path, section, config())
		assert.Error(t, err)
	})

	t.Run("validate fail", func(t *testing.T) {
		t.Parallel()

		cfg := config()
		cfg["host"] = "host.com/"
		_, err := newEdgegridConfig("", "", cfg)
		assert.Error(t, err)
	})
}
