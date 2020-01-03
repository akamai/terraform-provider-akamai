package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"fmt"
	"strconv"
)

//
// Handle Operations on gtm datacenters
// Based on 1.3 schema
//

// Datacenter represents a GTM datacenter
type Datacenter struct {
	City                       string      `json:"city,omitempty"`
	CloneOf                    int         `json:"cloneOf,omitempty"`
	CloudServerTargeting       bool        `json:"cloudServerTargeting,omitempty"`
	Continent                  string      `json:"continent,omitempty"`
	Country                    string      `json:"country,omitempty"`
	DefaultLoadObject          *LoadObject `json:"defaultLoadObject,omitempty"`
	Latitude                   float64     `json:"latitude,omitempty"`
	Links                      []*Link     `json:"links,omitempty"`
	Longitude                  float64     `json:"longitude,omitempty"`
	Nickname                   string      `json:"nickname,omitempty"`
	PingInterval               int         `json:"pingInterval,omitempty"`
	PingPacketSize             int         `json:"pingPacketSize,omitempty"`
	DatacenterId               int         `json:"datacenterId,omitempty"`
	ScorePenalty               int         `json:"scorePenalty,omitempty"`
	ServermonitorLivenessCount int         `json:"servermonitorLivenessCount,omitempty"`
	ServermonitorLoadCount     int         `json:"servermonitorLoadCount,omitempty"`
	ServermonitorPool          string      `json:"servermonitorPool,omitempty"`
	StateOrProvince            string      `json:"stateOrProvince,omitempty"`
	Virtual                    bool        `json:"virtual,omitempty"`
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

	printHttpResponse(res, true)

	// Network
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domainName,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

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

	printHttpResponse(res, true)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Datacenter",
			name:             strconv.Itoa(dc.DatacenterId),
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

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

	printHttpResponse(res, true)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Datacenter",
			name:             strconv.Itoa(dc.DatacenterId),
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

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
