package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"errors"
	"fmt"
	"strconv"
	"strings"
)

//
// Handle Operations on gtm datacenters
// Based on 1.4 schema
//

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
	DatacenterId                  int         `json:"datacenterId"`
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
func NewDatacenterResponse() *DatacenterResponse {
	dcResp := &DatacenterResponse{}
	return dcResp
}

// NewDatacenter creates a new Datacenter object
func NewDatacenter() *Datacenter {
	dc := &Datacenter{}
	return dc
}

// ListDatacenters retreieves all Datacenters
func ListDatacenters(domainName string) ([]*Datacenter, error) {
	dcs := &DatacenterList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Datacenter"}
	} else {
		err = client.BodyJSON(res, dcs)
		if err != nil {
			return nil, err
		}

		return dcs.DatacenterItems, nil
	}
}

// GetDatacenter retrieves a Datacenter with the given name. NOTE: Id arg is int!
func GetDatacenter(dcID int, domainName string) (*Datacenter, error) {

	dc := NewDatacenter()
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dcID)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpRequest(req, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Datacenter", name: strconv.Itoa(dcID)}
	} else {
		err = client.BodyJSON(res, dc)
		if err != nil {
			return nil, err
		}

		return dc, nil
	}
}

// Create the datacenter identified by the receiver argument in the specified domain.
func (dc *Datacenter) Create(domainName string) (*DatacenterResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters", domainName),
		dc,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)

	// Network
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domainName,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domainName, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := NewDatacenterResponse()
	// Unmarshall whole response body for updated DC and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil

}

var MapDefaultDC int = 5400
var Ipv4DefaultDC int = 5401
var Ipv6DefaultDC int = 5402

// Create Default Datacenter for Maps
func CreateMapsDefaultDatacenter(domainName string) (*Datacenter, error) {

	return createDefaultDC(MapDefaultDC, domainName)

}

// Create Default Datacenter for IPv4 Selector
func CreateIPv4DefaultDatacenter(domainName string) (*Datacenter, error) {

	return createDefaultDC(Ipv4DefaultDC, domainName)

}

// Create Default Datacenter for IPv6 Selector
func CreateIPv6DefaultDatacenter(domainName string) (*Datacenter, error) {

	return createDefaultDC(Ipv6DefaultDC, domainName)

}

// Worker function to create Default Datacenter identified id in the specified domain.
func createDefaultDC(defaultID int, domainName string) (*Datacenter, error) {

	if defaultID != MapDefaultDC && defaultID != Ipv4DefaultDC && defaultID != Ipv6DefaultDC {
		return nil, errors.New("Invalid default datacenter id provided for creation")
	}
	// check if already exists
	dc, err := GetDatacenter(defaultID, domainName)
	if err == nil {
		return dc, err
	} else {
		if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Datacenter") {
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
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		defaultURL,
		nil,
	)
	if err != nil {
		return nil, err
	}
	setVersionHeader(req, schemaVersion)
	printHttpRequest(req, true)
	res, err := client.Do(Config, req)
	// Network
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domainName,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}
	printHttpResponse(res, true)
	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domainName, apiErrorMessage: err.Detail, err: err}
	}
	responseBody := NewDatacenterResponse()
	// Unmarshall whole response body for updated DC and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Resource, nil

}

// Update the datacenter identified in the receiver argument in the provided domain.
func (dc *Datacenter) Update(domainName string) (*ResponseStatus, error) {

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterId)),
		dc,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Datacenter",
			name:             strconv.Itoa(dc.DatacenterId),
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Datacenter", name: string(dc.DatacenterId), apiErrorMessage: err.Detail, err: err}
	}

	responseBody := NewDatacenterResponse()
	// Unmarshall whole response body for updated entity and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil
}

// Delete the datacenter identified by the receiver argument from the domain specified.
func (dc *Datacenter) Delete(domainName string) (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s/datacenters/%s", domainName, strconv.Itoa(dc.DatacenterId)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Datacenter",
			name:             strconv.Itoa(dc.DatacenterId),
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Datacenter", name: string(dc.DatacenterId), apiErrorMessage: err.Detail, err: err}
	}

	responseBody := NewDatacenterResponse()
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil
}
