package papi

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Properties is a collection of PAPI Property resources
type Properties struct {
	client.Resource
	Properties struct {
		Items []*Property `json:"items"`
	} `json:"properties"`
}

// NewProperties creates a new Properties
func NewProperties() *Properties {
	properties := &Properties{}
	properties.Init()

	return properties
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (properties *Properties) PostUnmarshalJSON() error {
	properties.Init()

	for key, property := range properties.Properties.Items {
		properties.Properties.Items[key].parent = properties
		if err := property.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	properties.Complete <- true

	return nil
}

// GetProperties populates Properties with property data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listproperties
// Endpoint: GET /papi/v1/properties/{?contractId,groupId}
func (properties *Properties) GetProperties(contract *Contract, group *Group) error {
	if contract == nil {
		contract = NewContract(NewContracts())
		contract.ContractID = group.ContractIDs[0]
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties?groupId=%s&contractId=%s",
			group.GroupID,
			contract.ContractID,
		),
		nil,
	)

	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, properties); err != nil {
		return err
	}

	return nil
}

// AddProperty adds a property to the collection, if the property already exists
// in the collection it will be replaced.
func (properties *Properties) AddProperty(newProperty *Property) {
	if newProperty.PropertyID != "" {
		for key, property := range properties.Properties.Items {
			if property.PropertyID == newProperty.PropertyID {
				properties.Properties.Items[key] = newProperty
				return
			}
		}
	}

	newProperty.parent = properties

	properties.Properties.Items = append(properties.Properties.Items, newProperty)
}

// FindProperty finds a property by ID within the collection
func (properties *Properties) FindProperty(id string) (*Property, error) {
	var property *Property
	var propertyFound bool
	for _, property = range properties.Properties.Items {
		if property.PropertyID == id {
			propertyFound = true
			break
		}
	}

	if !propertyFound {
		return nil, fmt.Errorf("Unable to find property: \"%s\"", id)
	}

	return property, nil
}

// NewProperty creates a new property associated with the collection
func (properties *Properties) NewProperty(contract *Contract, group *Group) *Property {
	property := NewProperty(properties)

	properties.AddProperty(property)

	property.Contract = contract
	property.Group = group
	go property.Contract.GetContract()
	go property.Group.GetGroup()
	go (func(property *Property) {
		groupCompleted := <-property.Group.Complete
		contractCompleted := <-property.Contract.Complete
		property.Complete <- (groupCompleted && contractCompleted)
	})(property)

	return property
}

// Property represents a PAPI Property
type Property struct {
	client.Resource
	parent            *Properties
	AccountID         string             `json:"accountId,omitempty"`
	Contract          *Contract          `json:"-"`
	Group             *Group             `json:"-"`
	ContractID        string             `json:"contractId,omitempty"`
	GroupID           string             `json:"groupId,omitempty"`
	PropertyID        string             `json:"propertyId,omitempty"`
	PropertyName      string             `json:"propertyName"`
	LatestVersion     int                `json:"latestVersion,omitempty"`
	StagingVersion    int                `json:"stagingVersion,omitempty"`
	ProductionVersion int                `json:"productionVersion,omitempty"`
	Note              string             `json:"note,omitempty"`
	ProductID         string             `json:"productId,omitempty"`
	RuleFormat        string             `json:"ruleFormat",omitempty`
	CloneFrom         *ClonePropertyFrom `json:"cloneFrom"`
}

// NewProperty creates a new Property
func NewProperty(parent *Properties) *Property {
	property := &Property{parent: parent, Group: &Group{}, Contract: &Contract{}}
	property.Init()
	return property
}

// PreMarshalJSON is called before JSON marshaling
//
// See: jsonhooks-v1/json.Marshal()
func (property *Property) PreMarshalJSON() error {
	property.GroupID = property.Group.GroupID
	property.ContractID = property.Contract.ContractID
	return nil
}

// GetProperty populates a Property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaproperty
// Endpoint: GET /papi/v1/properties/{propertyId}{?contractId,groupId}
func (property *Property) GetProperty() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s",
			property.PropertyID,
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

	newProperties := NewProperties()
	if err := client.BodyJSON(res, newProperties); err != nil {
		return err
	}

	property.AccountID = newProperties.Properties.Items[0].AccountID
	property.Contract = newProperties.Properties.Items[0].Contract
	property.Group = newProperties.Properties.Items[0].Group
	property.ContractID = newProperties.Properties.Items[0].ContractID
	property.GroupID = newProperties.Properties.Items[0].GroupID
	property.PropertyID = newProperties.Properties.Items[0].PropertyID
	property.PropertyName = newProperties.Properties.Items[0].PropertyName
	property.LatestVersion = newProperties.Properties.Items[0].LatestVersion
	property.StagingVersion = newProperties.Properties.Items[0].StagingVersion
	property.ProductionVersion = newProperties.Properties.Items[0].ProductionVersion
	property.Note = newProperties.Properties.Items[0].Note
	property.ProductID = newProperties.Properties.Items[0].ProductID
	property.RuleFormat = newProperties.Properties.Items[0].RuleFormat
	property.CloneFrom = newProperties.Properties.Items[0].CloneFrom

	return nil
}

// GetActivations retrieves activation data for a given property
//
// See: Activations.GetActivations()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listactivations
// Endpoint: GET /papi/v1/properties/{propertyId}/activations/{?contractId,groupId}
func (property *Property) GetActivations() (*Activations, error) {
	activations := NewActivations()

	if err := activations.GetActivations(property); err != nil {
		return nil, err
	}

	return activations, nil
}

// GetAvailableBehaviors retrieves available behaviors for a given property
//
// See: AvailableBehaviors.GetAvailableBehaviors
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listavailablebehaviors
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/available-behaviors{?contractId,groupId}
func (property *Property) GetAvailableBehaviors() (*AvailableBehaviors, error) {
	behaviors := NewAvailableBehaviors()
	if err := behaviors.GetAvailableBehaviors(property); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// GetRules retrieves rules for a property
//
// See: Rules.GetRules
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletree
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (property *Property) GetRules() (*Rules, error) {
	rules := NewRules()

	if err := rules.GetRules(property); err != nil {
		return nil, err
	}

	return rules, nil
}

// GetRulesDigest fetches the Etag for a rule tree
//
// See: Rules.GetRulesDigest()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletreesdigest
// Endpoint: HEAD /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (property *Property) GetRulesDigest() (string, error) {
	rules := NewRules()
	return rules.GetRulesDigest(property)
}

// GetVersions retrieves all versions for a a given property
//
// See: Versions.GetVersions()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listversions
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{?contractId,groupId}
func (property *Property) GetVersions() (*Versions, error) {
	versions := NewVersions()
	err := versions.GetVersions(property)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// GetLatestVersion gets the latest active version, optionally of a given network
//
// See: Versions.GetLatestVersion()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getthelatestversion
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/latest{?contractId,groupId,activatedOn}
func (property *Property) GetLatestVersion(activatedOn NetworkValue) (*Version, error) {
	versions := NewVersions()
	versions.PropertyID = property.PropertyID

	return versions.GetLatestVersion(activatedOn)
}

// GetHostnames retrieves hostnames assigned to a given property
//
// If no version is given, the latest version is used
//
// See: Hostnames.GetHostnames()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getpropertyversionhostnames
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/hostnames/{?contractId,groupId}
func (property *Property) GetHostnames(version *Version) (*Hostnames, error) {
	hostnames := NewHostnames()
	hostnames.PropertyID = property.PropertyID
	hostnames.ContractID = property.Contract.ContractID
	hostnames.GroupID = property.Group.GroupID

	if version == nil {
		var err error
		version, err = property.GetLatestVersion("")
		if err != nil {
			return nil, err
		}
	}
	err := hostnames.GetHostnames(version)
	if err != nil {
		return nil, err
	}

	return hostnames, nil
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (property *Property) PostUnmarshalJSON() error {
	property.Init()

	property.Contract = NewContract(NewContracts())
	property.Contract.ContractID = property.ContractID

	property.Group = NewGroup(NewGroups())
	property.Group.GroupID = property.GroupID

	go property.Group.GetGroup()
	go property.Contract.GetContract()

	go (func(property *Property) {
		contractComplete := <-property.Contract.Complete
		groupComplete := <-property.Group.Complete
		property.Complete <- (contractComplete && groupComplete)
	})(property)

	return nil
}

// Save will create a property, optionally cloned from another property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#createorcloneaproperty
// Endpoint: POST /papi/v1/properties/{?contractId,groupId}
func (property *Property) Save() error {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/papi/v1/properties?contractId=%s&groupId=%s",
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		property,
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

	var location client.JSONBody
	if err = client.BodyJSON(res, &location); err != nil {
		return err
	}

	req, err = client.NewRequest(
		Config,
		"GET",
		location["propertyLink"].(string),
		nil,
	)
	if err != nil {
		return err
	}

	res, err = client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	properties := NewProperties()
	if err = client.BodyJSON(res, properties); err != nil {
		return err
	}

	property.AccountID = properties.Properties.Items[0].AccountID
	property.Contract = properties.Properties.Items[0].Contract
	property.Group = properties.Properties.Items[0].Group
	property.ContractID = properties.Properties.Items[0].ContractID
	property.GroupID = properties.Properties.Items[0].GroupID
	property.PropertyID = properties.Properties.Items[0].PropertyID
	property.PropertyName = properties.Properties.Items[0].PropertyName
	property.LatestVersion = properties.Properties.Items[0].LatestVersion
	property.StagingVersion = properties.Properties.Items[0].StagingVersion
	property.ProductionVersion = properties.Properties.Items[0].ProductionVersion
	property.Note = properties.Properties.Items[0].Note
	property.ProductID = properties.Properties.Items[0].ProductID
	property.CloneFrom = properties.Properties.Items[0].CloneFrom

	return nil
}

// Activate activates a given property
//
// If acknowledgeWarnings is true and warnings are returned on the first attempt,
// a second attempt is made, acknowledging the warnings.
//
// See: Activation.Save()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#activateaproperty
// Endpoint: POST /papi/v1/properties/{propertyId}/activations/{?contractId,groupId}
func (property *Property) Activate(activation *Activation, acknowledgeWarnings bool) error {
	return activation.Save(property, acknowledgeWarnings)
}

// Delete a property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#removeaproperty
// Endpoint: DELETE /papi/v1/properties/{propertyId}{?contractId,groupId}
func (property *Property) Delete() error {
	// /papi/v1/properties/{propertyId}{?contractId,groupId}
	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/papi/v1/properties/%s",
			property.PropertyID,
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

	return nil
}

// ClonePropertyFrom represents
type ClonePropertyFrom struct {
	client.Resource
	PropertyID           string `json:"propertyId"`
	Version              int    `json:"version"`
	CopyHostnames        bool   `json:"copyHostnames,omitempty"`
	CloneFromVersionEtag string `json:"cloneFromVersionEtag,omitempty"`
}

// NewClonePropertyFrom creates a new ClonePropertyFrom
func NewClonePropertyFrom() *ClonePropertyFrom {
	clonePropertyFrom := &ClonePropertyFrom{}
	clonePropertyFrom.Init()

	return clonePropertyFrom
}
