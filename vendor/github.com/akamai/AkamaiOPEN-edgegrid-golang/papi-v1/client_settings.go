package papi

import "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

// ClientSettings represents the PAPI client settings resource
type ClientSettings struct {
	client.Resource
	RuleFormat string `json:"ruleFormat"`
}

// NewClientSettings creates a new ClientSettings
func NewClientSettings() *ClientSettings {
	clientSettings := &ClientSettings{}
	clientSettings.Init()

	return clientSettings
}

// GetClientSettings populates ClientSettings
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getclientsettings
// Endpoint: GET /papi/v1/client-settings
func (clientSettings *ClientSettings) GetClientSettings() error {
	req, err := client.NewRequest(Config, "GET", "/papi/v1/client-settings", nil)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err := client.BodyJSON(res, clientSettings); err != nil {
		return err
	}

	return nil
}

// Save updates client settings
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#updateclientsettings
// Endpoint: PUT /papi/v1/client-settings
func (clientSettings *ClientSettings) Save() error {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		"/papi/v1/client-settings",
		clientSettings,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	newClientSettings := NewClientSettings()
	if err := client.BodyJSON(res, newClientSettings); err != nil {
		return err
	}

	clientSettings.RuleFormat = newClientSettings.RuleFormat

	return nil
}
