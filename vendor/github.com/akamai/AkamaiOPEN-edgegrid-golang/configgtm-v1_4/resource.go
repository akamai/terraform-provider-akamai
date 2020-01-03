package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"fmt"
)

//
// Handle Operations on gtm resources
// Based on 1.4 schema
//

// ResourceInstance
type ResourceInstance struct {
	DatacenterId         int  `json:"datacenterId"`
	UseDefaultLoadObject bool `json:"useDefaultLoadObject,omitempty"`
	LoadObject
}

// Resource represents a GTM resource
type Resource struct {
	Type                        string              `json:"type"`
	HostHeader                  string              `json:"hostHeader,omitempty"`
	LeastSquaresDecay           int                 `json:"leastSquaresDecay,omitempty"`
	Description                 string              `json:"description,omitempty"`
	LeaderString                string              `json:"leaderString,omitempty"`
	ConstrainedProperty         string              `json:"constrainedProperty,omitempty"`
	ResourceInstances           []*ResourceInstance `json:"resourceInstances,omitempty"`
	AggregationType             string              `json:"aggregationType"`
	Links                       []*Link             `json:"links,omitempty"`
	LoadImbalancePercentage     float64             `json:"loadImbalancePercentage,omitempty"`
	UpperBound                  int                 `json:"upperBound,omitempty"`
	Name                        string              `json:"name"`
	MaxUMultiplicativeIncrement float64             `json:"maxUMultiplicativeIncrement,omitempty"`
	DecayRate                   float64             `json:"decayRate,omitempty"`
}

// ResourceList is the structure returned by List Resources
type ResourceList struct {
	ResourceItems []*Resource `json:"items"`
}

// NewResourceInstance instantiates a new ResourceInstance.
func (rsrc *Resource) NewResourceInstance(dcID int) *ResourceInstance {

	return &ResourceInstance{DatacenterId: dcID}

}

// NewResource creates a new Resource object.
func NewResource(name string) *Resource {
	resource := &Resource{Name: name}
	return resource
}

// ListResources retreieves all Resources in the specified domain.
func ListResources(domainName string) ([]*Resource, error) {
	rsrcs := &ResourceList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/resources", domainName),
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
		return nil, CommonError{entityName: "Resources"}
	}
	err = client.BodyJSON(res, rsrcs)
	if err != nil {
		return nil, err
	}

	return rsrcs.ResourceItems, nil

}

// GetResource retrieves a Resource with the given name in the specified domain.
func GetResource(name, domainName string) (*Resource, error) {
	rsc := NewResource(name)
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, name),
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
		return nil, CommonError{entityName: "Resource", name: name}
	} else {
		err = client.BodyJSON(res, rsc)
		if err != nil {
			return nil, err
		}

		return rsc, nil
	}
}

// Create the resource identified by the receiver argument in the specified domain.
func (rsrc *Resource) Create(domainName string) (*ResourceResponse, error) {

	// Use common code. Any specific validation needed?

	return rsrc.save(domainName)

}

// Update the resourceidentified in the receiver argument in the specified domain.
func (rsrc *Resource) Update(domainName string) (*ResponseStatus, error) {

	// common code

	stat, err := rsrc.save(domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save Resource in given domain. Common path for Create and Update.
func (rsrc *Resource) save(domainName string) (*ResourceResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name),
		rsrc,
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
			entityName:       "Resource",
			name:             rsrc.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Resource", name: rsrc.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResourceResponse{}
	// Unmarshall whole response body for updated entity and in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil

}

// Delete the resource identified in the receiver argument from the specified domain.
func (rsrc *Resource) Delete(domainName string) (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name),
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
			entityName:       "Resource",
			name:             rsrc.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Resource", name: rsrc.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil

}
