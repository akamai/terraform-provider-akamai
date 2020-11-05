package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// ClientSettings contains operations available on ClientSettings resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#clientsettingsgroup
	ClientSettings interface {
		// GetClientSettings returns client's settings.
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getclientsettings
		GetClientSettings(context.Context) (*ClientSettingsBody, error)
		// UpdateClientSettings updates client's settings.
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#putclientsettings
		UpdateClientSettings(context.Context, ClientSettingsBody) (*ClientSettingsBody, error)
	}

	// ClientSettingsBody represents both the request and response bodies for operating on client settings resource
	ClientSettingsBody struct {
		RuleFormat  string `json:"ruleFormat"`
		UsePrefixes bool   `json:"usePrefixes"`
	}
)

var (
	ErrGetClientSettings    = errors.New("fetching client settings")
	ErrUpdateClientSettings = errors.New("updating client settings")
)

// GetClientSettings is used to list the client settings
func (p *papi) GetClientSettings(ctx context.Context) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("GetClientSettings")

	getURL := "/papi/v1/client-settings"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetClientSettings, err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetClientSettings, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetClientSettings, p.Error(resp))
	}

	return &clientSettings, nil
}

// UpdateClientSettings is used to update the client settings
// fixme body structure
func (p *papi) UpdateClientSettings(ctx context.Context, params ClientSettingsBody) (*ClientSettingsBody, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateClientSettings")

	putURL := "/papi/v1/client-settings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateClientSettings, err)
	}

	var clientSettings ClientSettingsBody
	resp, err := p.Exec(req, &clientSettings, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateClientSettings, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetClientSettings, p.Error(resp))
	}

	return &clientSettings, nil
}
