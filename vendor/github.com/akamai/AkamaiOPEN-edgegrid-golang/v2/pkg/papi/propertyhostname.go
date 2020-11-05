package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PropertyVersionHostnames contains operations available on PropertyVersionHostnames resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#propertyversionhostnamesgroup
	PropertyVersionHostnames interface {
		// GetPropertyVersionHostnames lists all the hostnames assigned to a property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyversionhostnames
		GetPropertyVersionHostnames(context.Context, GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error)

		// CreatePropertyVersionHostnames lists all the hostnames assigned to a property version
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#putpropertyversionhostnames
		UpdatePropertyVersionHostnames(context.Context, UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error)
	}

	// GetPropertyVersionHostnamesRequest contains parameters required to list property version hostnames
	GetPropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
	}

	// GetPropertyVersionHostnamesResponse contains all property version hostnames associated to the given parameters
	GetPropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// HostnameResponseItems contains the response body for GetPropertyVersionHostnamesResponse
	HostnameResponseItems struct {
		Items []Hostname `json:"items"`
	}

	// Hostname contains information about each of the HostnameResponseItems
	Hostname struct {
		CnameType      HostnameCnameType `json:"cnameType"`
		EdgeHostnameID string            `json:"edgeHostnameId"`
		CnameFrom      string            `json:"cnameFrom"`
		CnameTo        string            `json:"cnameTo"`
	}

	// UpdatePropertyVersionHostnamesRequest contains parameters required to update the set of hostname entries for a property version
	UpdatePropertyVersionHostnamesRequest struct {
		PropertyID        string
		PropertyVersion   int
		ContractID        string
		GroupID           string
		ValidateHostnames bool
		Hostnames         []Hostname
	}

	// UpdatePropertyVersionHostnamesResponse contains information about each of the HostnameRequestItems
	UpdatePropertyVersionHostnamesResponse struct {
		AccountID       string                `json:"accountId"`
		ContractID      string                `json:"contractId"`
		GroupID         string                `json:"groupId"`
		PropertyID      string                `json:"propertyId"`
		PropertyVersion int                   `json:"propertyVersion"`
		Etag            string                `json:"etag"`
		Hostnames       HostnameResponseItems `json:"hostnames"`
	}

	// HostnameCnameType represents HostnameCnameType enum
	HostnameCnameType string
)

const (
	// HostnameCnameTypeEdgeHostname const
	HostnameCnameTypeEdgeHostname HostnameCnameType = "EDGE_HOSTNAME"
)

// Validate validates GetPropertyVersionHostnamesRequest
func (ph GetPropertyVersionHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ph.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ph.PropertyVersion, validation.Required),
	}.Filter()
}

// Validate validates UpdatePropertyVersionHostnamesRequest
func (ch UpdatePropertyVersionHostnamesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(ch.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(ch.PropertyVersion, validation.Required),
		"Hostnames":       validation.Validate(ch.Hostnames, validation.Required),
	}.Filter()
}

var (
	ErrGetPropertyVersionHostnames    = errors.New("fetching hostnames")
	ErrUpdatePropertyVersionHostnames = errors.New("updating hostnames")
)

func (p *papi) GetPropertyVersionHostnames(ctx context.Context, params GetPropertyVersionHostnamesRequest) (*GetPropertyVersionHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetPropertyVersionHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetPropertyVersionHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/hostnames?contractId=%s&groupId=%s&validateHostnames=%t",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
		params.ValidateHostnames)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPropertyVersionHostnames, err)
	}

	var hostnames GetPropertyVersionHostnamesResponse
	resp, err := p.Exec(req, &hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPropertyVersionHostnames, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}

func (p *papi) UpdatePropertyVersionHostnames(ctx context.Context, params UpdatePropertyVersionHostnamesRequest) (*UpdatePropertyVersionHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdatePropertyVersionHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("UpdatePropertyVersionHostnames")

	putURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%v/hostnames?contractId=%s&groupId=%s&validateHostnames=%t",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
		params.ValidateHostnames,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdatePropertyVersionHostnames, err)
	}

	var hostnames UpdatePropertyVersionHostnamesResponse
	resp, err := p.Exec(req, &hostnames, params.Hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdatePropertyVersionHostnames, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdatePropertyVersionHostnames, p.Error(resp))
	}

	return &hostnames, nil
}
