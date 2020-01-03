package client

import (
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	config := edgegrid.Config{
		Host:         "https://httpbin.org",
		AccessToken:  "local-config",
		ClientSecret: "local-config",
		ClientToken:  "local-config",
	}

	req, err := NewRequest(config, "GET", "/headers", nil)
	assert.NotNil(t, req)
	assert.NoError(t, err)
	assert.Equal(t, "https://httpbin.org/headers", req.URL.String())
	verifyResponseConfig(t, config, req)
}

func TestNewRequestWithSwitchedAccount(t *testing.T) {
	config := edgegrid.Config{
		Host:         "https://httpbin.org",
		AccessToken:  "local-config",
		ClientSecret: "local-config",
		ClientToken:  "local-config",
		AccountKey:   "ABC-DEF",
	}

	req, err := NewRequest(config, "GET", "/headers", nil)
	assert.NotNil(t, req)
	assert.NoError(t, err)
	assert.Equal(t, "https://httpbin.org/headers?accountSwitchKey=ABC-DEF", req.URL.String())
	verifyResponseConfig(t, config, req)
}

func TestNewJSONRequest(t *testing.T) {
	config := edgegrid.Config{
		Host:         "https://httpbin.org",
		AccessToken:  "local-config",
		ClientSecret: "local-config",
		ClientToken:  "local-config",
	}

	req, err := NewJSONRequest(config, "GET", "/headers", config)
	assert.NotNil(t, req)
	assert.NoError(t, err)
	assert.Equal(t, "https://httpbin.org/headers", req.URL.String())
	verifyResponseConfig(t, config, req)
}

func TestNewJSONRequestWithSwitchedAccount(t *testing.T) {
	config := edgegrid.Config{
		Host:         "https://httpbin.org",
		AccessToken:  "local-config",
		ClientSecret: "local-config",
		ClientToken:  "local-config",
		AccountKey:   "ABC-DEF",
	}

	req, err := NewJSONRequest(config, "GET", "/headers", config)
	assert.NotNil(t, req)
	assert.NoError(t, err)
	assert.Equal(t, "https://httpbin.org/headers?accountSwitchKey=ABC-DEF", req.URL.String())
	verifyResponseConfig(t, config, req)
}

func verifyResponseConfig(t *testing.T, config edgegrid.Config, req *http.Request) {
	resp, err := Do(config, req)
	assert.NotNil(t, resp)
	assert.NoError(t, err)

	json := make(map[string]interface{})
	BodyJSON(resp, &json)

	assert.True(t, strings.Contains(json["headers"].(map[string]interface{})["Authorization"].(string), "local-config"))
}
