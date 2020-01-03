package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"fmt"
)

//
// Handle Operations on gtm geomaps
// Based on 1.4 schema
//

// GeoAssigment represents a GTM geo assignment element
type GeoAssignment struct {
	DatacenterBase
	Countries []string `json:"countries,omitempty"`
}

// GeoMap  represents a GTM GeoMap
type GeoMap struct {
	DefaultDatacenter *DatacenterBase  `json:"defaultDatacenter"`
	Assignments       []*GeoAssignment `json:"assignments,omitempty"`
	Name              string           `json:"name"`
	Links             []*Link          `json:"links,omitempty"`
}

// GeoMapList represents the returned GTM GeoMap List body
type GeoMapList struct {
	GeoMapItems []*GeoMap `json:"items"`
}

// NewGeoMap creates a new GeoMap object
func NewGeoMap(name string) *GeoMap {
	geomap := &GeoMap{Name: name}
	return geomap
}

// ListGeoMap retreieves all GeoMaps
func ListGeoMaps(domainName string) ([]*GeoMap, error) {
	geos := &GeoMapList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps", domainName),
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
		return nil, CommonError{entityName: "geoMap"}
	}
	err = client.BodyJSON(res, geos)
	if err != nil {
		return nil, err
	}

	return geos.GeoMapItems, nil

}

// GetGeoMap retrieves a GeoMap with the given name.
func GetGeoMap(name, domainName string) (*GeoMap, error) {
	geo := NewGeoMap(name)

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, name),
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
		return nil, CommonError{entityName: "GeographicMap", name: name}
	} else {
		err = client.BodyJSON(res, geo)
		if err != nil {
			return nil, err
		}

		return geo, nil
	}
}

// Instantiate new Assignment struct
func (geo *GeoMap) NewAssignment(dcID int, nickname string) *GeoAssignment {
	geoAssign := &GeoAssignment{}
	geoAssign.DatacenterId = dcID
	geoAssign.Nickname = nickname

	return geoAssign

}

// Instantiate new Default Datacenter Struct
func (geo *GeoMap) NewDefaultDatacenter(dcID int) *DatacenterBase {
	return &DatacenterBase{DatacenterId: dcID}
}

// Create GeoMap in provided domain
func (geo *GeoMap) Create(domainName string) (*GeoMapResponse, error) {

	// Use common code. Any specific validation needed?

	return geo.save(domainName)

}

// Update GeoMap in given domain
func (geo *GeoMap) Update(domainName string) (*ResponseStatus, error) {

	// common code

	stat, err := geo.save(domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save GeoMap in given domain. Common path for Create and Update.
func (geo *GeoMap) save(domainName string) (*GeoMapResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, geo.Name),
		geo,
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
			entityName:       "geographicMap",
			name:             geo.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "geographicMap", name: geo.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &GeoMapResponse{}
	// Unmarshall whole response body for updated entity and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// Delete GeoMap method
func (geo *GeoMap) Delete(domainName string) (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s/geographic-maps/%s", domainName, geo.Name),
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
			entityName:       "geographicMap",
			name:             geo.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "geographicMap", name: geo.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil
}
