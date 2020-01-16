package configgtm

import (
	"fmt"
	"net/http"
)

//
// Common data types and methods
// Based on 1.3 schemas
//

// Append url args to req
func appendReqArgs(req *http.Request, queryArgs map[string]string) {

	// Look for optional args
	if len(queryArgs) > 0 {
		q := req.URL.Query()
		for argName, argVal := range queryArgs {
			q.Add(argName, argVal)
		}
		req.URL.RawQuery = q.Encode()
	}

}

// default schema version
// TODO: retrieve from environment or elsewhere in Service Init
var schemaVersion string = "1.3"

// internal method to set version. passed in as string
func setVersionHeader(req *http.Request, version string) {

	req.Header.Set("Accept", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))

	if req.Method != "GET" {
		req.Header.Set("Content-Type", fmt.Sprintf("application/vnd.config-gtm.v%s+json", version))
	}

	return

}

// response Status is returned on Create, Update or Delete operations for all entity types
type ResponseStatus struct {
	ChangeId              string  `json:"changeId"`
	Links                 *[]Link `json:"links"`
	Message               string  `json:"message"`
	PassingValidation     bool    `json:"passingValidation"`
	PropagationStatus     string  `json:"propagationStatus"`
	PropagationStatusDate string  `json:"propagationStatusDate"`
}

// NewResponseStatus returns a new ResponseStatus struct
func NewResponseStatus() *ResponseStatus {

	return &ResponseStatus{}

}

// Generic response structs
type ResponseBody struct {
	Resource interface{}     `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// Response structs by Entity Type
type DomainResponse struct {
	Resource *Domain         `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

type DatacenterResponse struct {
	Status   *ResponseStatus `json:"status"`
	Resource *Datacenter     `json:"resource"`
}

type PropertyResponse struct {
	Resource *Property       `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

type ResourceResponse struct {
	Resource *Resource       `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

type CidrMapResponse struct {
	Resource *CidrMap        `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

type GeoMapResponse struct {
	Resource *GeoMap         `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

type AsMapResponse struct {
	Resource *AsMap          `json:"resource"`
	Status   *ResponseStatus `json:"status"`
}

// Probably THE most common type
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

//
type LoadObject struct {
	LoadObject     string   `json:"loadObject, omitempty"`
	LoadObjectPort int      `json:"loadObjectPort, omitempty"`
	LoadServers    []string `json:"loadServers, omitempty"`
}

// NewLoadObject returns a new LoadObject structure
func NewLoadObject() *LoadObject {
	return &LoadObject{}
}

type DatacenterBase struct {
	Nickname     string `json:"nickname"`
	DatacenterId int    `json:"datacenterId"`
}

// NewDatacenterBase returns a new DatacenterBase structure
func NewDatacenterBase() *DatacenterBase {
	return &DatacenterBase{}
}
