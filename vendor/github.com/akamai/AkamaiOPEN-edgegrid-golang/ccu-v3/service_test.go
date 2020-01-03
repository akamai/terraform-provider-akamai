package ccu

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/stretchr/testify/assert"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        false,
	}
)

func TestInit(t *testing.T) {
	Init(config)

	assert.Equal(t, config.Host, Config.Host)
	assert.Equal(t, config.AccessToken, Config.AccessToken)
	assert.Equal(t, config.ClientToken, Config.ClientToken)
	assert.Equal(t, config.ClientSecret, Config.ClientSecret)
	assert.Equal(t, config.MaxBody, Config.MaxBody)
	assert.Equal(t, config.Debug, Config.Debug)
}
