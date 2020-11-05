package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm cidrmaps
// Based on 1.4 schema
//

// CidrMaps contains operations available on a Cidrmap resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type CidrMaps interface {
	// NewCidrMap creates a new CidrMap object.
	NewCidrMap(context.Context, string) *CidrMap
	// Instantiate new Assignment struct
	NewCidrAssignment(context.Context, *CidrMap, int, string) *CidrAssignment
	// ListCidrMaps retreieves all CidrMaps
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getcidrmaps
	ListCidrMaps(context.Context, string) ([]*CidrMap, error)
	// GetCidrMap retrieves a CidrMap with the given name.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getcidrmap
	GetCidrMap(context.Context, string, string) (*CidrMap, error)
	// Create the datacenter identified by the receiver argument in the specified domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putcidrmap
	CreateCidrMap(context.Context, *CidrMap, string) (*CidrMapResponse, error)
	// Delete the datacenter identified by the receiver argument from the domain specified.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#deletecidrmap
	DeleteCidrMap(context.Context, *CidrMap, string) (*ResponseStatus, error)
	// Update the datacenter identified in the receiver argument in the provided domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putcidrmap
	UpdateCidrMap(context.Context, *CidrMap, string) (*ResponseStatus, error)
}

//CidrAssignment represents a GTM cidr assignment element
type CidrAssignment struct {
	DatacenterBase
	Blocks []string `json:"blocks"`
}

// CidrMap  represents a GTM cidrMap element
type CidrMap struct {
	DefaultDatacenter *DatacenterBase   `json:"defaultDatacenter"`
	Assignments       []*CidrAssignment `json:"assignments,omitempty"`
	Name              string            `json:"name"`
	Links             []*Link           `json:"links,omitempty"`
}

// CidrMapList represents a GTM returned cidrmap list body
type CidrMapList struct {
	CidrMapItems []*CidrMap `json:"items"`
}

// Validate validates CidrMap
func (cid *CidrMap) Validate() error {
	if len(cid.Name) < 1 {
		return fmt.Errorf("CidrMap is missing Name")
	}
	if cid.DefaultDatacenter == nil {
		return fmt.Errorf("CidrMap is missing DefaultDatacenter")
	}

	return nil
}

// NewCidrMap creates a new CidrMap object
func (p *gtm) NewCidrMap(ctx context.Context, name string) *CidrMap {

	logger := p.Log(ctx)
	logger.Debug("NewCidrMap")

	cidrmap := &CidrMap{Name: name}
	return cidrmap
}

// ListCidrMap retreieves all CidrMaps
func (p *gtm) ListCidrMaps(ctx context.Context, domainName string) ([]*CidrMap, error) {

	logger := p.Log(ctx)
	logger.Debug("ListCidrMaps")

	var cidrs CidrMapList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListCidrMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &cidrs)
	if err != nil {
		return nil, fmt.Errorf("ListCidrMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return cidrs.CidrMapItems, nil
}

// GetCidrMap retrieves a CidrMap with the given name.
func (p *gtm) GetCidrMap(ctx context.Context, name, domainName string) (*CidrMap, error) {

	logger := p.Log(ctx)
	logger.Debug("GetCidrMap")

	var cidr CidrMap
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetCidrMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &cidr)
	if err != nil {
		return nil, fmt.Errorf("GetCidrMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &cidr, nil
}

// Instantiate new Assignment struct
func (p *gtm) NewCidrAssignment(ctx context.Context, cidr *CidrMap, dcid int, nickname string) *CidrAssignment {

	logger := p.Log(ctx)
	logger.Debug("NewCidrAssignment")

	cidrAssign := &CidrAssignment{}
	cidrAssign.DatacenterId = dcid
	cidrAssign.Nickname = nickname

	return cidrAssign
}

// Create CidrMap in provided domain
func (p *gtm) CreateCidrMap(ctx context.Context, cidr *CidrMap, domainName string) (*CidrMapResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateCidrMap")

	// Use common code. Any specific validation needed?
	return cidr.save(ctx, p, domainName)
}

// Update CidrMap in given domain
func (p *gtm) UpdateCidrMap(ctx context.Context, cidr *CidrMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateCidrMap")

	// common code
	stat, err := cidr.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save CidrMap in given domain. Common path for Create and Update.
func (cidr *CidrMap) save(ctx context.Context, p *gtm, domainName string) (*CidrMapResponse, error) {

	if err := cidr.Validate(); err != nil {
		return nil, fmt.Errorf("CidrMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, cidr.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create AsMap request: %w", err)
	}

	var mapresp CidrMapResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp, cidr)
	if err != nil {
		return nil, fmt.Errorf("CidrMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &mapresp, nil
}

// Delete CidrMap method
func (p *gtm) DeleteCidrMap(ctx context.Context, cidr *CidrMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteCidrMap")

	if err := cidr.Validate(); err != nil {
		logger.Errorf("CidrMap validation failed. %w", err)
		return nil, fmt.Errorf("CidrMap validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/cidr-maps/%s", domainName, cidr.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var mapresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp)
	if err != nil {
		return nil, fmt.Errorf("CidrMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return mapresp.Status, nil
}
