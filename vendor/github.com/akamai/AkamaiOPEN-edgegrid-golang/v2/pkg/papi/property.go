package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Properties contains operations available on Property resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#propertiesgroup
	Properties interface {
		GetProperties(ctx context.Context, r GetPropertiesRequest) (*GetPropertiesResponse, error)
		CreateProperty(ctx context.Context, params CreatePropertyRequest) (*CreatePropertyResponse, error)
		GetProperty(ctx context.Context, params GetPropertyRequest) (*GetPropertyResponse, error)
		RemoveProperty(ctx context.Context, params RemovePropertyRequest) (*RemovePropertyResponse, error)
	}

	// PropertyCloneFrom optionally identifies another property instance to clone when making a POST request to create a new property
	PropertyCloneFrom struct {
		CloneFromVersionEtag string `json:"cloneFromVersionEtag,omitempty"`
		CopyHostnames        bool   `json:"copyHostnames,omitempty"`
		PropertyID           string `json:"propertyId"`
		Version              int    `json:"version"`
	}

	// Property contains configuration data to apply to edge content.
	Property struct {
		AccountID         string `json:"accountId"`
		AssetID           string `json:"assetId"`
		ContractID        string `json:"contractId"`
		GroupID           string `json:"groupId"`
		LatestVersion     int    `json:"latestVersion"`
		Note              string `json:"note"`
		ProductID         string `json:"productId"`
		ProductionVersion *int   `json:"productionVersion,omitempty"`
		PropertyID        string `json:"propertyId"`
		PropertyName      string `json:"propertyName"`
		RuleFormat        string `json:"ruleFormat"`
		StagingVersion    *int   `json:"stagingVersion,omitempty"`
	}

	// PropertiesItems is an array of properties
	PropertiesItems struct {
		Items []*Property `json:"items"`
	}

	// GetPropertiesRequest is the argument for GetProperties
	GetPropertiesRequest struct {
		ContractID string
		GroupID    string
	}

	// GetPropertiesResponse is the response for GetProperties
	GetPropertiesResponse struct {
		Properties PropertiesItems `json:"properties"`
	}

	// CreatePropertyRequest is passed to CreateProperty
	CreatePropertyRequest struct {
		ContractID string
		GroupID    string
		Property   PropertyCreate
	}

	// PropertyCreate represents a POST /property request body
	PropertyCreate struct {
		CloneFrom    *PropertyCloneFrom `json:"cloneFrom,omitempty"`
		ProductID    string             `json:"productId"`
		PropertyName string             `json:"propertyName"`
		RuleFormat   string             `json:"ruleFormat,omitempty"`
	}

	// CreatePropertyResponse is returned by CreateProperty
	CreatePropertyResponse struct {
		Response
		PropertyID   string
		PropertyLink string `json:"propertyLink"`
	}

	// GetPropertyRequest is the argument for GetProperty
	GetPropertyRequest struct {
		ContractID string
		GroupID    string
		PropertyID string
	}

	// GetPropertyResponse is the response for GetProperty
	GetPropertyResponse struct {
		Response
		Properties PropertiesItems `json:"properties"`
		Property   *Property       `json:"-"`
	}

	// RemovePropertyRequest is the argument for RemoveProperty
	RemovePropertyRequest struct {
		PropertyID string
		ContractID string
		GroupID    string
	}

	// RemovePropertyResponse is the response for GetProperties
	RemovePropertyResponse struct {
		Message string `json:"message"`
	}
)

// Validate validates GetPropertiesRequest
func (v GetPropertiesRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(v.ContractID, validation.Required),
		"GroupID":    validation.Validate(v.GroupID, validation.Required),
	}.Filter()
}

// Validate validates CreatePropertyRequest
func (v CreatePropertyRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(v.ContractID, validation.Required),
		"GroupID":    validation.Validate(v.GroupID, validation.Required),
		"Property":   validation.Validate(v.Property),
	}.Filter()
}

// Validate validates PropertyCreate
func (p PropertyCreate) Validate() error {
	return validation.Errors{
		"ProductID":    validation.Validate(p.ProductID, validation.Required),
		"PropertyName": validation.Validate(p.PropertyName, validation.Required),
		"CloneFrom":    validation.Validate(p.CloneFrom),
	}.Filter()
}

// Validate validates PropertyCloneFrom
func (c PropertyCloneFrom) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(c.PropertyID),
		"Version":    validation.Validate(c.Version),
	}.Filter()
}

// Validate validates GetPropertyRequest
func (v GetPropertyRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(v.PropertyID, validation.Required),
	}.Filter()
}

// Validate validates RemovePropertyRequest
func (v RemovePropertyRequest) Validate() error {
	return validation.Errors{
		"PropertyID": validation.Validate(v.PropertyID, validation.Required),
	}.Filter()
}

var (
	ErrGetProperties  = errors.New("fetching properties")
	ErrGetProperty    = errors.New("fetching property")
	ErrCreateProperty = errors.New("creating property")
	ErrRemoveProperty = errors.New("removing property")
)

func (p *papi) GetProperties(ctx context.Context, params GetPropertiesRequest) (*GetPropertiesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProperties, ErrStructValidation, err)
	}

	var rval GetPropertiesResponse

	logger := p.Log(ctx)
	logger.Debug("GetProperties")

	uri := fmt.Sprintf(
		"/papi/v1/properties?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProperties, err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProperties, p.Error(resp))
	}

	return &rval, nil
}

func (p *papi) CreateProperty(ctx context.Context, params CreatePropertyRequest) (*CreatePropertyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateProperty, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreateProperty")

	uri := fmt.Sprintf(
		"/papi/v1/properties?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateProperty, err)
	}

	var rval CreatePropertyResponse

	resp, err := p.Exec(req, &rval, params.Property)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateProperty, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateProperty, p.Error(resp))
	}

	id, err := ResponseLinkParse(rval.PropertyLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateProperty, ErrInvalidResponseLink, err)
	}
	rval.PropertyID = id

	return &rval, nil
}

func (p *papi) GetProperty(ctx context.Context, params GetPropertyRequest) (*GetPropertyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetProperty, ErrStructValidation, err)
	}

	var rval GetPropertyResponse

	logger := p.Log(ctx)
	logger.Debug("GetProperty")

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s",
		params.PropertyID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetProperty, err)
	}
	q := uri.Query()
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetProperty, err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetProperty, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetProperty, p.Error(resp))
	}

	if len(rval.Properties.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: PropertyID: %s", ErrGetProperty, ErrNotFound, params.PropertyID)
	}
	rval.Property = rval.Properties.Items[0]

	return &rval, nil
}

func (p *papi) RemoveProperty(ctx context.Context, params RemovePropertyRequest) (*RemovePropertyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrRemoveProperty, ErrStructValidation, err)
	}

	var rval RemovePropertyResponse

	logger := p.Log(ctx)
	logger.Debug("RemoveProperty")

	uri, err := url.Parse(fmt.Sprintf(
		"/papi/v1/properties/%s",
		params.PropertyID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed parse url: %s", ErrRemoveProperty, err)
	}
	q := uri.Query()
	if params.GroupID != "" {
		q.Add("groupId", params.GroupID)
	}
	if params.ContractID != "" {
		q.Add("contractId", params.ContractID)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRemoveProperty, err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRemoveProperty, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrRemoveProperty, p.Error(resp))
	}

	return &rval, nil
}
