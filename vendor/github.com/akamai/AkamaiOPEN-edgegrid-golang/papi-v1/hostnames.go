package papi

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Hostnames is a collection of Property Hostnames
type Hostnames struct {
	client.Resource
	AccountID       string `json:"accountId"`
	ContractID      string `json:"contractId"`
	GroupID         string `json:"groupId"`
	PropertyID      string `json:"propertyId"`
	PropertyVersion int    `json:"propertyVersion"`
	Etag            string `json:"etag"`
	Hostnames       struct {
		Items []*Hostname `json:"items"`
	} `json:"hostnames"`
}

// NewHostnames creates a new Hostnames
func NewHostnames() *Hostnames {
	hostnames := &Hostnames{}
	hostnames.Init()

	return hostnames
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (hostnames *Hostnames) PostUnmarshalJSON() error {
	hostnames.Init()

	for key, hostname := range hostnames.Hostnames.Items {
		hostnames.Hostnames.Items[key].parent = hostnames
		if err := hostname.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	hostnames.Complete <- true

	return nil
}

// GetHostnames retrieves hostnames assigned to a given property
//
// If no version is given, the latest version is used
//
// See: Property.GetHostnames()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listapropertyshostnames
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/hostnames/{?contractId,groupId}
func (hostnames *Hostnames) GetHostnames(version *Version) error {
	if version == nil {
		property := NewProperty(NewProperties())
		property.PropertyID = hostnames.PropertyID
		err := property.GetProperty()
		if err != nil {
			return err
		}

		version, err = property.GetLatestVersion("")
		if err != nil {
			return err
		}
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/hostnames/?contractId=%s&groupId=%s",
			hostnames.PropertyID,
			version.PropertyVersion,
			hostnames.ContractID,
			hostnames.GroupID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, hostnames); err != nil {
		return err
	}

	return nil
}

// NewHostname creates a new Hostname within a given Hostnames
func (hostnames *Hostnames) NewHostname() *Hostname {
	hostname := NewHostname(hostnames)
	hostnames.Hostnames.Items = append(hostnames.Hostnames.Items, hostname)
	return hostname
}

// Save updates a properties hostnames
func (hostnames *Hostnames) Save() error {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/hostnames?contractId=%s&groupId=%s",
			hostnames.PropertyID,
			hostnames.PropertyVersion,
			hostnames.ContractID,
			hostnames.GroupID,
		),
		hostnames.Hostnames.Items,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, hostnames); err != nil {
		return err
	}

	return nil
}

// Hostname represents a property hostname resource
type Hostname struct {
	client.Resource
	parent         *Hostnames
	CnameType      CnameTypeValue `json:"cnameType"`
	EdgeHostnameID string         `json:"edgeHostnameId"`
	CnameFrom      string         `json:"cnameFrom"`
	CnameTo        string         `json:"cnameTo,omitempty"`
}

// NewHostname creates a new Hostname
func NewHostname(parent *Hostnames) *Hostname {
	hostname := &Hostname{parent: parent, CnameType: CnameTypeEdgeHostname}
	hostname.Init()

	return hostname
}

// CnameTypeValue is used to create an "enum" of possible Hostname.CnameType values
type CnameTypeValue string

const (
	// CnameTypeEdgeHostname Hostname.CnameType value EDGE_HOSTNAME
	CnameTypeEdgeHostname CnameTypeValue = "EDGE_HOSTNAME"
)
