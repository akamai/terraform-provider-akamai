package papi

import (
	"fmt"
	"io/ioutil"

        edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/xeipuuv/gojsonschema"
)

// AvailableCriteria represents a collection of available rule criteria
type AvailableCriteria struct {
	client.Resource
	ContractID        string `json:"contractId"`
	GroupID           string `json:"groupId"`
	ProductID         string `json:"productId"`
	RuleFormat        string `json:"ruleFormat"`
	AvailableCriteria struct {
		Items []struct {
			Name       string `json:"name"`
			SchemaLink string `json:"schemaLink"`
		} `json:"items"`
	} `json:"availableCriteria"`
}

// NewAvailableCriteria creates a new AvailableCriteria
func NewAvailableCriteria() *AvailableCriteria {
	availableCriteria := &AvailableCriteria{}
	availableCriteria.Init()

	return availableCriteria
}

// GetAvailableCriteria retrieves criteria available for a given property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listavailablecriteria
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/available-criteria{?contractId,groupId}
func (availableCriteria *AvailableCriteria) GetAvailableCriteria(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/available-criteria?contractId=%s&groupId=%s",
			property.PropertyID,
			property.LatestVersion,
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, availableCriteria); err != nil {
		return err
	}

	return nil
}

// AvailableBehaviors represents a collection of available rule behaviors
type AvailableBehaviors struct {
	client.Resource
	ContractID string `json:"contractId"`
	GroupID    string `json:"groupId"`
	ProductID  string `json:"productId"`
	RuleFormat string `json:"ruleFormat"`
	Behaviors  struct {
		Items []AvailableBehavior `json:"items"`
	} `json:"behaviors"`
}

// NewAvailableBehaviors creates a new AvailableBehaviors
func NewAvailableBehaviors() *AvailableBehaviors {
	availableBehaviors := &AvailableBehaviors{}
	availableBehaviors.Init()

	return availableBehaviors
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (availableBehaviors *AvailableBehaviors) PostUnmarshalJSON() error {
	availableBehaviors.Init()

	for key := range availableBehaviors.Behaviors.Items {
		availableBehaviors.Behaviors.Items[key].parent = availableBehaviors
	}

	availableBehaviors.Complete <- true

	return nil
}

// GetAvailableBehaviors retrieves available behaviors for a given property
//
// See: Property.GetAvailableBehaviors
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listavailablebehaviors
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/available-behaviors{?contractId,groupId}
func (availableBehaviors *AvailableBehaviors) GetAvailableBehaviors(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/available-behaviors?contractId=%s&groupId=%s",
			property.PropertyID,
			property.LatestVersion,
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, availableBehaviors); err != nil {
		return err
	}

	return nil
}

// AvailableBehavior represents an available behavior resource
type AvailableBehavior struct {
	client.Resource
	parent     *AvailableBehaviors
	Name       string `json:"name"`
	SchemaLink string `json:"schemaLink"`
}

// NewAvailableBehavior creates a new AvailableBehavior
func NewAvailableBehavior(parent *AvailableBehaviors) *AvailableBehavior {
	availableBehavior := &AvailableBehavior{parent: parent}
	availableBehavior.Init()

	return availableBehavior
}

// GetSchema retrieves the JSON schema for an available behavior
func (behavior *AvailableBehavior) GetSchema() (*gojsonschema.Schema, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		behavior.SchemaLink,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	schemaBytes, _ := ioutil.ReadAll(res.Body)
	schemaBody := string(schemaBytes)
	loader := gojsonschema.NewStringLoader(schemaBody)
	schema, err := gojsonschema.NewSchema(loader)

	return schema, err
}
