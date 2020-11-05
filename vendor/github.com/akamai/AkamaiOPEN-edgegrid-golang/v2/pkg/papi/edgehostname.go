package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// EdgeHostnames contains operations available on EdgeHostnames resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#edgehostnamesgroup
	EdgeHostnames interface {
		// GetEdgeHostnames fetches a list of edge hostnames
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getedgehostnames
		GetEdgeHostnames(context.Context, GetEdgeHostnamesRequest) (*GetEdgeHostnamesResponse, error)

		// GetEdgeHostname fetches edge hostname with given ID
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getedgehostname
		GetEdgeHostname(context.Context, GetEdgeHostnameRequest) (*GetEdgeHostnamesResponse, error)

		// CreateEdgeHostname creates a new edge hostname
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#postedgehostnames
		CreateEdgeHostname(context.Context, CreateEdgeHostnameRequest) (*CreateEdgeHostnameResponse, error)
	}

	// GetEdgeHostnamesRequest contains query params used for listing edge hostnames
	GetEdgeHostnamesRequest struct {
		ContractID string
		GroupID    string
		Options    []string
	}

	// GetEdgeHostnameRequest contains path and query params used to fetch specific edge hostname
	GetEdgeHostnameRequest struct {
		EdgeHostnameID string
		ContractID     string
		GroupID        string
		Options        []string
	}

	// GetEdgeHostnamesResponse contains data received by calling GetEdgeHostnames or GetEdgeHostname
	GetEdgeHostnamesResponse struct {
		AccountID     string            `json:"accountId"`
		ContractID    string            `json:"contractId"`
		GroupID       string            `json:"groupId"`
		EdgeHostnames EdgeHostnameItems `json:"edgeHostnames"`
		EdgeHostname  EdgeHostnameGetItem
	}

	// EdgeHostnameItems contains a list of EdgeHostnames
	EdgeHostnameItems struct {
		Items []EdgeHostnameGetItem `json:"items"`
	}

	// EdgeHostnameGetItem contains GET details for edge hostname
	EdgeHostnameGetItem struct {
		ID                string    `json:"edgeHostnameId"`
		Domain            string    `json:"edgeHostnameDomain"`
		ProductID         string    `json:"productId"`
		DomainPrefix      string    `json:"domainPrefix"`
		DomainSuffix      string    `json:"domainSuffix"`
		Status            string    `json:"status,omitempty"`
		Secure            bool      `json:"secure"`
		IPVersionBehavior string    `json:"ipVersionBehavior"`
		UseCases          []UseCase `json:"useCases,omitempty"`
	}

	// UseCase contains UseCase data
	UseCase struct {
		Option  string `json:"option"`
		Type    string `json:"type"`
		UseCase string `json:"useCase"`
	}

	// CreateEdgeHostnameRequest contains query params and body required for creation of new edge hostname
	CreateEdgeHostnameRequest struct {
		ContractID   string
		GroupID      string
		Options      []string
		EdgeHostname EdgeHostnameCreate
	}

	// EdgeHostnameCreate contains body of edge hostname POST request
	EdgeHostnameCreate struct {
		ProductID         string    `json:"productId"`
		DomainPrefix      string    `json:"domainPrefix"`
		DomainSuffix      string    `json:"domainSuffix"`
		Secure            bool      `json:"secure,omitempty"`
		SecureNetwork     string    `json:"secureNetwork,omitempty"`
		SlotNumber        int       `json:"slotNumber,omitempty"`
		IPVersionBehavior string    `json:"ipVersionBehavior"`
		CertEnrollmentID  int       `json:"certEnrollmentId,omitempty"`
		UseCases          []UseCase `json:"useCases,omitempty"`
	}

	// CreateEdgeHostnameResponse contains a link returned after creating new edge hostname and DI of this hostname
	CreateEdgeHostnameResponse struct {
		EdgeHostnameLink string `json:"edgeHostnameLink"`
		EdgeHostnameID   string `json:"-"`
	}
)

const (
	// EHSecureNetworkStandardTLS constant
	EHSecureNetworkStandardTLS = "STANDARD_TLS"
	// EHSecureNetworkSharedCert constant
	EHSecureNetworkSharedCert = "SHARED_CERT"
	// EHSecureNetworkEnhancedTLS constant
	EHSecureNetworkEnhancedTLS = "ENHANCED_TLS"

	// EHIPVersionV4 constant
	EHIPVersionV4 = "IPV4"
	// EHIPVersionV6Performance constant
	EHIPVersionV6Performance = "IPV6_PERFORMANCE"
	// EHIPVersionV6Compliance constant
	EHIPVersionV6Compliance = "IPV6_COMPLIANCE"

	// UseCaseGlobal constant
	UseCaseGlobal = "GLOBAL"
)

// Validate validates CreateEdgeHostnameRequest
func (eh CreateEdgeHostnameRequest) Validate() error {
	return validation.Errors{
		"ContractID":   validation.Validate(eh.ContractID, validation.Required),
		"GroupID":      validation.Validate(eh.GroupID, validation.Required),
		"EdgeHostname": validation.Validate(eh.EdgeHostname),
	}.Filter()
}

// Validate validates EdgeHostnameCreate
func (eh EdgeHostnameCreate) Validate() error {
	return validation.Errors{
		"DomainPrefix": validation.Validate(eh.DomainPrefix, validation.Required),
		"DomainSuffix": validation.Validate(eh.DomainSuffix, validation.Required,
			validation.When(eh.SecureNetwork == EHSecureNetworkStandardTLS, validation.In("edgesuite.net")),
			validation.When(eh.SecureNetwork == EHSecureNetworkSharedCert, validation.In("akamaized.net")),
			validation.When(eh.SecureNetwork == EHSecureNetworkEnhancedTLS, validation.In("edgekey.net")),
		),
		"ProductID":         validation.Validate(eh.ProductID, validation.Required),
		"CertEnrollmentID":  validation.Validate(eh.CertEnrollmentID, validation.Required.When(eh.SecureNetwork == EHSecureNetworkEnhancedTLS)),
		"IPVersionBehavior": validation.Validate(eh.IPVersionBehavior, validation.Required, validation.In(EHIPVersionV4, EHIPVersionV6Performance, EHIPVersionV6Compliance)),
		"SecureNetwork":     validation.Validate(eh.SecureNetwork, validation.In(EHSecureNetworkStandardTLS, EHSecureNetworkSharedCert, EHSecureNetworkEnhancedTLS)),
		"UseCases":          validation.Validate(eh.UseCases),
	}.Filter()
}

// Validate validates UseCase
func (uc UseCase) Validate() error {
	return validation.Errors{
		"Option":  validation.Validate(uc.Option, validation.Required),
		"Type":    validation.Validate(uc.Type, validation.Required, validation.In(UseCaseGlobal)),
		"UseCase": validation.Validate(uc.UseCase, validation.Required),
	}.Filter()
}

// Validate validates GetEdgeHostnamesRequest
func (eh GetEdgeHostnamesRequest) Validate() error {
	return validation.Errors{
		"ContractID": validation.Validate(eh.ContractID, validation.Required),
		"GroupID":    validation.Validate(eh.GroupID, validation.Required),
	}.Filter()
}

// Validate validates GetEdgeHostnameRequest
func (eh GetEdgeHostnameRequest) Validate() error {
	return validation.Errors{
		"EdgeHostnameID": validation.Validate(eh.EdgeHostnameID, validation.Required),
		"ContractID":     validation.Validate(eh.ContractID, validation.Required),
		"GroupID":        validation.Validate(eh.GroupID, validation.Required),
	}.Filter()
}

var (
	ErrGetEdgeHostnames   = errors.New("fetching edge hostnames")
	ErrGetEdgeHostname    = errors.New("fetching edge hostname")
	ErrCreateEdgeHostname = errors.New("creating edge hostname")
)

// GetEdgeHostnames id used to list edge hostnames for provided group and contract IDs
func (p *papi) GetEdgeHostnames(ctx context.Context, params GetEdgeHostnamesRequest) (*GetEdgeHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeHostnames, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetEdgeHostnames")

	getURL := fmt.Sprintf(
		"/papi/v1/edgehostnames?contractId=%s&groupId=%s",
		params.ContractID,
		params.GroupID,
	)
	if len(params.Options) > 0 {
		getURL = fmt.Sprintf("%s&options=%s", getURL, strings.Join(params.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeHostnames, err)
	}

	var edgeHostnames GetEdgeHostnamesResponse
	resp, err := p.Exec(req, &edgeHostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeHostnames, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetEdgeHostnames, p.Error(resp))
	}

	return &edgeHostnames, nil
}

// GetEdgeHostname id used to fetch edge hostname with given ID for provided group and contract IDs
func (p *papi) GetEdgeHostname(ctx context.Context, params GetEdgeHostnameRequest) (*GetEdgeHostnamesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetEdgeHostname, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetEdgeHostname")

	getURL := fmt.Sprintf(
		"/papi/v1/edgehostnames/%s?contractId=%s&groupId=%s",
		params.EdgeHostnameID,
		params.ContractID,
		params.GroupID,
	)
	if len(params.Options) > 0 {
		getURL = fmt.Sprintf("%s&options=%s", getURL, strings.Join(params.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetEdgeHostname, err)
	}

	var edgeHostname GetEdgeHostnamesResponse
	resp, err := p.Exec(req, &edgeHostname)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetEdgeHostname, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}
	if len(edgeHostname.EdgeHostnames.Items) == 0 {
		return nil, fmt.Errorf("%s: %w: EdgeHostnameID: %s", ErrGetEdgeHostname, ErrNotFound, params.EdgeHostnameID)
	}
	edgeHostname.EdgeHostname = edgeHostname.EdgeHostnames.Items[0]

	return &edgeHostname, nil
}

// CreateEdgeHostname id used to create new edge hostname for provided group and contract IDs
func (p *papi) CreateEdgeHostname(ctx context.Context, r CreateEdgeHostnameRequest) (*CreateEdgeHostnameResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEdgeHostname, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("CreateEdgeHostname")

	createURL := fmt.Sprintf(
		"/papi/v1/edgehostnames?contractId=%s&groupId=%s",
		r.ContractID,
		r.GroupID,
	)
	if len(r.Options) > 0 {
		createURL = fmt.Sprintf("%s&options=%s", createURL, strings.Join(r.Options, ","))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateEdgeHostname, err)
	}

	var createResponse CreateEdgeHostnameResponse
	resp, err := p.Exec(req, &createResponse, r.EdgeHostname)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateEdgeHostname, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateEdgeHostname, p.Error(resp))
	}
	id, err := ResponseLinkParse(createResponse.EdgeHostnameLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateEdgeHostname, ErrInvalidResponseLink, err)
	}
	createResponse.EdgeHostnameID = id
	return &createResponse, nil
}
