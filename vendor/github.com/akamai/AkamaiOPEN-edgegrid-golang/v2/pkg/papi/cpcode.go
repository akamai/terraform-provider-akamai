package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CPCodes contains operations available on CPCode resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#cpcodesgroup
	CPCodes interface {
		// GetCPCodes lists all available CP codes
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcpcodes
		GetCPCodes(context.Context, GetCPCodesRequest) (*GetCPCodesResponse, error)

		// GetCPCode gets CP code with provided ID
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcpcode
		GetCPCode(context.Context, GetCPCodeRequest) (*GetCPCodesResponse, error)

		// CreateCPCode creates a new CP code
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postcpcodes
		CreateCPCode(context.Context, CreateCPCodeRequest) (*CreateCPCodeResponse, error)
	}

	// CPCode contains CP code resource data
	CPCode struct {
		ID          string   `json:"cpcodeId"`
		Name        string   `json:"cpcodeName"`
		CreatedDate string   `json:"createdDate"`
		ProductIDs  []string `json:"productIds"`
	}

	// CPCodeItems contains a list of CPCode items
	CPCodeItems struct {
		Items []CPCode `json:"items"`
	}

	// GetCPCodesResponse is a response returned while fetching CP codes
	GetCPCodesResponse struct {
		AccountID  string      `json:"accountId"`
		ContractID string      `json:"contractId"`
		GroupID    string      `json:"groupId"`
		CPCodes    CPCodeItems `json:"cpcodes"`
		CPCode     CPCode
	}

	// CreateCPCodeRequest contains data required to create CP code (both request body and group/contract infromation
	CreateCPCodeRequest struct {
		ContractID string
		GroupID    string
		CPCode     CreateCPCode
	}

	// CreateCPCode contains the request body for CP code creation
	CreateCPCode struct {
		ProductID  string `json:"productId"`
		CPCodeName string `json:"cpcodeName"`
	}

	// CreateCPCodeResponse contains the response from CP code creation as well as the ID of created resource
	CreateCPCodeResponse struct {
		CPCodeLink string `json:"cpcodeLink"`
		CPCodeID   string `json:"-"`
	}

	// GetCPCodeRequest gets details about a CP code.
	GetCPCodeRequest struct {
		CPCodeID   string
		ContractID string
		GroupID    string
	}

	// GetCPCodesRequest contains parameters require to list/create CP codes
	// GroupID and ContractID are required as part of every CP code operation, ID is required only for operating on specific CP code
	GetCPCodesRequest struct {
		ContractID string
		GroupID    string
	}
)

// Validate validates GetCPCodesRequest
func (cp GetCPCodesRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
	}.Filter()
}

// Validate validates GetCPCodeRequest
func (cp GetCPCodeRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
		"CPCodeID":   validation.Validate(cp.CPCodeID, validation.Required),
	}.Filter()
}

// Validate validates CreateCPCodeRequest
func (cp CreateCPCodeRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(cp.ContractID, validation.Required),
		"GroupID":    validation.Validate(cp.GroupID, validation.Required),
		"CPCode":     validation.Validate(cp.CPCode, validation.Required),
	}.Filter()
}

// Validate validates CreateCPCode
func (cp CreateCPCode) Validate() error {
	return validation.Errors{
		"ProductID":  validation.Validate(cp.ProductID, validation.Required),
		"CPCodeName": validation.Validate(cp.CPCodeName, validation.Required),
	}.Filter()
}

var (
	ErrGetCPCodes   = errors.New("fetching CP Codes")
	ErrGetCPCode    = errors.New("fetching CP Code")
	ErrCreateCPCode = errors.New("creating CP Code")
)

// GetCPCodes is used to list all available CP codes for given group and contract
func (p *papi) GetCPCodes(ctx context.Context, params GetCPCodesRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCodes, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCodes")

	getURL := fmt.Sprintf(
		"/papi/v1/cpcodes?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCodes, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCodes, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCodes, p.Error(resp))
	}

	return &cpCodes, nil
}

// GetCPCodes is used to fetch a CP code with provided ID
func (p *papi) GetCPCode(ctx context.Context, params GetCPCodeRequest) (*GetCPCodesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCPCode, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetCPCode")

	getURL := fmt.Sprintf("/papi/v1/cpcodes/%s?contractId=%s&groupId=%s", params.CPCodeID, params.ContractID, params.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCPCode, err)
	}

	var cpCodes GetCPCodesResponse
	resp, err := p.Exec(req, &cpCodes)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCPCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetCPCode, p.Error(resp))
	}
	if len(cpCodes.CPCodes.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: CPCodeID: %s", ErrGetCPCode, ErrNotFound, params.CPCodeID)
	}
	cpCodes.CPCode = cpCodes.CPCodes.Items[0]

	return &cpCodes, nil
}

// CreateCPCode creates a new CP code with provided CreateCPCodeRequest data
func (p *papi) CreateCPCode(ctx context.Context, r CreateCPCodeRequest) (*CreateCPCodeResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", ErrCreateCPCode, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreateCPCode")

	createURL := fmt.Sprintf("/papi/v1/cpcodes?contractId=%s&groupId=%s", r.ContractID, r.GroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCPCode, err)
	}

	var createResponse CreateCPCodeResponse
	resp, err := p.Exec(req, &createResponse, r.CPCode)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCPCode, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateCPCode, p.Error(resp))
	}
	id, err := ResponseLinkParse(createResponse.CPCodeLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCPCode, ErrInvalidResponseLink, err)
	}
	createResponse.CPCodeID = id
	return &createResponse, nil
}
