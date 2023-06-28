package akamai

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEdgegridConfig(t *testing.T) {
	t.Parallel()

	path := "testdata/edgerc"
	section := "default"
	config := func() *edgegrid.Config {
		return &edgegrid.Config{
			Host:         "host.com",
			AccessToken:  "access_token",
			ClientToken:  "client_token",
			ClientSecret: "client_secret",
			MaxBody:      0,
			AccountKey:   "",
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
		cfg.Host = "host.com/"
		_, err := newEdgegridConfig("", "", cfg)
		assert.Error(t, err)
	})
}
