package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"fmt"
)

//
// Handle Operations on gtm asmaps
// Based on 1.3 schema
//

// AsAssignment represents a GTM asmap assignment structure
type AsAssignment struct {
	DatacenterBase
	AsNumbers []int64 `json:"asNumbers"`
}

// AsMap  represents a GTM AsMap
type AsMap struct {
	DefaultDatacenter *DatacenterBase `json:"defaultDatacenter"`
	Assignments       []*AsAssignment `json:"assignments,omitempty"`
	Name              string          `json:"name,omitempty"`
	Links             []*Link         `json:"links,omitempty"`
}

// NewAsMap creates a new asMap
func NewAsMap(name string) *AsMap {
	asmap := &AsMap{Name: name}
	return asmap
}

// GetAsMap retrieves a asMap with the given name.
func GetAsMap(name, domainName string) (*AsMap, error) {
	as := NewAsMap(name)
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, name),
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
		return nil, CommonError{entityName: "asMap", name: name}
	} else {
		err = client.BodyJSON(res, as)
		if err != nil {
			return nil, err
		}

		return as, nil
	}
}

// Instantiate new Assignment struct
func (as *AsMap) NewAssignment(dcID int, nickname string) *AsAssignment {
	asAssign := &AsAssignment{}
	asAssign.DatacenterId = dcID
	asAssign.Nickname = nickname

	return asAssign

}

// Instantiate new Default Datacenter Struct
func (as *AsMap) NewDefaultDatacenter(dcID int) *DatacenterBase {
	return &DatacenterBase{DatacenterId: dcID}
}

// Create asMap in provided domain
func (as *AsMap) Create(domainName string) (*AsMapResponse, error) {

	// Use common code. Any specific validation needed?

	return as.save(domainName)

}

// Update AsMap in given domain
func (as *AsMap) Update(domainName string) (*ResponseStatus, error) {

	// common code

	stat, err := as.save(domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save AsMap in given domain. Common path for Create and Update.
func (as *AsMap) save(domainName string) (*AsMapResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, as.Name),
		as,
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
			entityName:       "asMap",
			name:             as.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "asMap", name: as.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &AsMapResponse{}
	// Unmarshall whole response body for updated entity and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// Delete AsMap method
func (as *AsMap) Delete(domainName string) (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s/as-maps/%s", domainName, as.Name),
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
			entityName:       "asMap",
			name:             as.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "asMap", name: as.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil
}
