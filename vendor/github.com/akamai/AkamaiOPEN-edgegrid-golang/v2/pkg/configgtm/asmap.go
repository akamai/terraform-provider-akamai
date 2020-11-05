package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm asmaps
// Based on 1.4 schema
//

// ASMaps contains operations available on a ASmap resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type ASMaps interface {
	// NewAsMap creates a new AsMap object.
	NewAsMap(context.Context, string) *AsMap
	// Instantiate new Assignment struct
	NewASAssignment(context.Context, *AsMap, int, string) *AsAssignment
	// ListAsMaps retreieves all AsMaps
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getasmaps
	ListAsMaps(context.Context, string) ([]*AsMap, error)
	// GetAsMap retrieves a AsMap with the given name.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getasmap
	GetAsMap(context.Context, string, string) (*AsMap, error)
	// Create the datacenter identified by the receiver argument in the specified domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putasmap
	CreateAsMap(context.Context, *AsMap, string) (*AsMapResponse, error)
	// Delete the datacenter identified by the receiver argument from the domain specified.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#deleteasmap
	DeleteAsMap(context.Context, *AsMap, string) (*ResponseStatus, error)
	// Update the datacenter identified in the receiver argument in the provided domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putasmap
	UpdateAsMap(context.Context, *AsMap, string) (*ResponseStatus, error)
}

// AsAssignment represents a GTM asmap assignment structure
type AsAssignment struct {
	DatacenterBase
	AsNumbers []int64 `json:"asNumbers"`
}

// AsMap  represents a GTM AsMap
type AsMap struct {
	DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
	Assignments       []*AsAssignment `json:"assignments,omitempty"`
	Name              string          `json:"name"`
	Links             []*Link         `json:"links,omitempty"`
}

// AsMapList represents the returned GTM AsMap List body
type AsMapList struct {
	AsMapItems []*AsMap `json:"items"`
}

// Validate validates AsMap
func (asm *AsMap) Validate() error {

	if len(asm.Name) < 1 {
		return fmt.Errorf("AsMap is missing Name")
	}
	if asm.DefaultDatacenter == nil {
		return fmt.Errorf("AsMap is missing DefaultDatacenter")
	}

	return nil
}

// NewAsMap creates a new asMap
func (p *gtm) NewAsMap(ctx context.Context, name string) *AsMap {

	logger := p.Log(ctx)
	logger.Debug("NewAsMap")

	asmap := &AsMap{Name: name}
	return asmap
}

// ListAsMaps retreieves all AsMaps
func (p *gtm) ListAsMaps(ctx context.Context, domainName string) ([]*AsMap, error) {

	logger := p.Log(ctx)
	logger.Debug("ListAsMaps")

	var aslist AsMapList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListAsMaps request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &aslist)
	if err != nil {
		return nil, fmt.Errorf("ListAsMaps request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return aslist.AsMapItems, nil
}

// GetAsMap retrieves a asMap with the given name.
func (p *gtm) GetAsMap(ctx context.Context, name, domainName string) (*AsMap, error) {

	logger := p.Log(ctx)
	logger.Debug("GetAsMap")

	var as AsMap
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAsMap request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &as)
	if err != nil {
		return nil, fmt.Errorf("GetAsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &as, nil
}

// Instantiate new Assignment struct
func (p *gtm) NewASAssignment(ctx context.Context, as *AsMap, dcID int, nickname string) *AsAssignment {

	logger := p.Log(ctx)
	logger.Debug("NewAssignment")

	asAssign := &AsAssignment{}
	asAssign.DatacenterId = dcID
	asAssign.Nickname = nickname

	return asAssign
}

// Create asMap in provided domain
func (p *gtm) CreateAsMap(ctx context.Context, as *AsMap, domainName string) (*AsMapResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateAsMap")

	// Use common code. Any specific validation needed?
	return as.save(ctx, p, domainName)
}

// Update AsMap in given domain
func (p *gtm) UpdateAsMap(ctx context.Context, as *AsMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateAsMap")

	// common code
	stat, err := as.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Save AsMap in given domain. Common path for Create and Update.
func (as *AsMap) save(ctx context.Context, p *gtm, domainName string) (*AsMapResponse, error) {

	if err := as.Validate(); err != nil {
		return nil, fmt.Errorf("AsMap validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, as.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create AsMap request: %w", err)
	}

	var mapresp AsMapResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp, as)
	if err != nil {
		return nil, fmt.Errorf("AsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &mapresp, nil
}

// Delete AsMap method
func (p *gtm) DeleteAsMap(ctx context.Context, as *AsMap, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteAsMap")

	if err := as.Validate(); err != nil {
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, as.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var mapresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &mapresp)
	if err != nil {
		return nil, fmt.Errorf("AsMap request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return mapresp.Status, nil
}
