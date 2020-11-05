package gtm

import (
	"context"
	"fmt"
	"net/http"

	"strconv"
)

//
// Handle Operations on gtm datacenters
// Based on 1.4 schema
//

// Datacenters contains operations available on a Datacenter resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type Datacenters interface {
	// NewDatacenterResponse instantiates a new DatacenterResponse structure
	NewDatacenterResponse(context.Context) *DatacenterResponse
	// NewDatacenter creates a new Datacenter object.
	NewDatacenter(context.Context) *Datacenter
	// ListDatacenters retreieves all Datacenters
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getdatacenters
	ListDatacenters(context.Context, string) ([]*Datacenter, error)
	// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getdatacenter
	GetDatacenter(context.Context, int, string) (*Datacenter, error)
	// Create the datacenter identified by the receiver argument in the specified domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#postdatacenters
	CreateDatacenter(context.Context, *Datacenter, string) (*DatacenterResponse, error)
	// Delete the datacenter identified by the receiver argument from the domain specified.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#deletedatacenter
	DeleteDatacenter(context.Context, *Datacenter, string) (*ResponseStatus, error)
	// Update the datacenter identified in the receiver argument in the provided domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putdatacenter
	UpdateDatacenter(context.Context, *Datacenter, string) (*ResponseStatus, error)
	// Create Default Datacenter for Maps
	CreateMapsDefaultDatacenter(context.Context, string) (*Datacenter, error)
	// Create Default Datacenter for IPv4 Selector
	CreateIPv4DefaultDatacenter(context.Context, string) (*Datacenter, error)
	// Create Default Datacenter for IPv6 Selector
	CreateIPv6DefaultDatacenter(context.Context, string) (*Datacenter, error)
}

// Datacenter represents a GTM datacenter
type Datacenter struct {
	City                          string      `json:"city,omitempty"`
	CloneOf                       int         `json:"cloneOf,omitempty"`
	CloudServerHostHeaderOverride bool        `json:"cloudServerHostHeaderOverride"`
	CloudServerTargeting          bool        `json:"cloudServerTargeting"`
	Continent                     string      `json:"continent,omitempty"`
	Country                       string      `json:"country,omitempty"`
	DefaultLoadObject             *LoadObject `json:"defaultLoadObject,omitempty"`
	Latitude                      float64     `json:"latitude,omitempty"`
	Links                         []*Link     `json:"links,omitempty"`
	Longitude                     float64     `json:"longitude,omitempty"`
	Nickname                      string      `json:"nickname,omitempty"`
	PingInterval                  int         `json:"pingInterval,omitempty"`
	PingPacketSize                int         `json:"pingPacketSize,omitempty"`
	DatacenterId                  int         `json:"datacenterId,omitempty"`
	ScorePenalty                  int         `json:"scorePenalty,omitempty"`
	ServermonitorLivenessCount    int         `json:"servermonitorLivenessCount,omitempty"`
	ServermonitorLoadCount        int         `json:"servermonitorLoadCount,omitempty"`
	ServermonitorPool             string      `json:"servermonitorPool,omitempty"`
	StateOrProvince               string      `json:"stateOrProvince,omitempty"`
	Virtual                       bool        `json:"virtual"`
}

type DatacenterList struct {
	DatacenterItems []*Datacenter `json:"items"`
}

// NewDatacenterResponse instantiates a new DatacenterResponse structure
func (p *gtm) NewDatacenterResponse(ctx context.Context) *DatacenterResponse {

	logger := p.Log(ctx)
	logger.Debug("NewDatacenterResponse")

	dcResp := &DatacenterResponse{}
	return dcResp
}

// NewDatacenter creates a new Datacenter object
func (p *gtm) NewDatacenter(ctx context.Context) *Datacenter {

	logger := p.Log(ctx)
	logger.Debug("NewDatacenter")

	dc := &Datacenter{}
	return dc
}

// ListDatacenters retreieves all Datacenters
func (p *gtm) ListDatacenters(ctx context.Context, domainName string) ([]*Datacenter, error) {

	logger := p.Log(ctx)
	logger.Debug("ListDatacenters")

	var dcs DatacenterList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListDatacenters request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dcs)
	if err != nil {
		return nil, fmt.Errorf("ListDatacenters request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return dcs.DatacenterItems, nil
}

// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
func (p *gtm) GetDatacenter(ctx context.Context, dcID int, domainName string) (*Datacenter, error) {

	logger := p.Log(ctx)
	logger.Debug("GetDatacenter")

	var dc Datacenter
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dcID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDatacenter request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dc)
	if err != nil {
		return nil, fmt.Errorf("GetDatacenter request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &dc, nil
}

// Create the datacenter identified by the receiver argument in the specified domain.
func (p *gtm) CreateDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*DatacenterResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateDatacenter")

	postURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Datacenter request: %w", err)
	}

	var dcresp DatacenterResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dcresp, dc)
	if err != nil {
		return nil, fmt.Errorf("Datacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &dcresp, nil
}

var MapDefaultDC int = 5400
var Ipv4DefaultDC int = 5401
var Ipv6DefaultDC int = 5402

// Create Default Datacenter for Maps
func (p *gtm) CreateMapsDefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateMapsDefaultDatacenter")

	return createDefaultDC(ctx, p, MapDefaultDC, domainName)
}

// Create Default Datacenter for IPv4 Selector
func (p *gtm) CreateIPv4DefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateIPv4DefaultDatacenter")

	return createDefaultDC(ctx, p, Ipv4DefaultDC, domainName)
}

// Create Default Datacenter for IPv6 Selector
func (p *gtm) CreateIPv6DefaultDatacenter(ctx context.Context, domainName string) (*Datacenter, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateIPv6DefaultDatacenter")

	return createDefaultDC(ctx, p, Ipv6DefaultDC, domainName)
}

// Worker function to create Default Datacenter identified id in the specified domain.
func createDefaultDC(ctx context.Context, p *gtm, defaultID int, domainName string) (*Datacenter, error) {

	if defaultID != MapDefaultDC && defaultID != Ipv4DefaultDC && defaultID != Ipv6DefaultDC {
		return nil, fmt.Errorf("Invalid default datacenter id provided for creation")
	}
	// check if already exists
	dc, err := p.GetDatacenter(ctx, defaultID, domainName)
	if err == nil {
		return dc, err
	} else {
		apiError, ok := err.(*Error)
		//if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Datacenter") {
		if !ok || apiError.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}
	defaultURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/", domainName)
	switch defaultID {
	case MapDefaultDC:
		defaultURL += "default-datacenter-for-maps"
	case Ipv4DefaultDC:
		defaultURL += "datacenter-for-ip-version-selector-ipv4"
	case Ipv6DefaultDC:
		defaultURL += "datacenter-for-ip-version-selector-ipv6"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Default Datacenter request: %w", err)
	}

	var dcresp DatacenterResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dcresp, "")
	if err != nil {
		return nil, fmt.Errorf("Default Datacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return dcresp.Resource, nil

}

// Update the datacenter identified in the receiver argument in the provided domain.
func (p *gtm) UpdateDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateDatacenter")

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterId))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Update Datacenter request: %w", err)
	}

	var dcresp DatacenterResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dcresp, dc)
	if err != nil {
		return nil, fmt.Errorf("Datacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return dcresp.Status, nil
}

// Delete the datacenter identified by the receiver argument from the domain specified.
func (p *gtm) DeleteDatacenter(ctx context.Context, dc *Datacenter, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteDatacenter")

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterId))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete Datacenter request: %w", err)
	}

	var dcresp DatacenterResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &dcresp)
	if err != nil {
		return nil, fmt.Errorf("Datacenter request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return dcresp.Status, nil
}
