package papi

import (
	"errors"
	"fmt"
        edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"time"
)

// Versions contains a collection of Property Versions
type Versions struct {
	client.Resource
	PropertyID   string `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	AccountID    string `json:"accountId"`
	ContractID   string `json:"contractId"`
	GroupID      string `json:"groupId"`
	Versions     struct {
		Items []*Version `json:"items"`
	} `json:"versions"`
	RuleFormat string `json:"ruleFormat,omitempty"`
}

// NewVersions creates a new Versions
func NewVersions() *Versions {
	version := &Versions{}
	version.Init()

	return version
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (versions *Versions) PostUnmarshalJSON() error {
	versions.Init()

	for key := range versions.Versions.Items {
		versions.Versions.Items[key].parent = versions
	}
	versions.Complete <- true

	return nil
}

// AddVersion adds or replaces a version within the collection
func (versions *Versions) AddVersion(version *Version) {
	if version.PropertyVersion != 0 {
		for key, v := range versions.Versions.Items {
			if v.PropertyVersion == version.PropertyVersion {
				versions.Versions.Items[key] = version
				return
			}
		}
	}

	versions.Versions.Items = append(versions.Versions.Items, version)
}

// GetVersions retrieves all versions for a a given property
//
// See: Property.GetVersions()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listversions
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{?contractId,groupId}
func (versions *Versions) GetVersions(property *Property) error {
	if property == nil {
		return errors.New("You must provide a property")
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions",
			property.PropertyID,
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

	if err = client.BodyJSON(res, versions); err != nil {
		return err
	}

	return nil
}

// GetLatestVersion retrieves the latest Version for a property
//
// See: Property.GetLatestVersion()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getthelatestversion
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/latest{?contractId,groupId,activatedOn}
func (versions *Versions) GetLatestVersion(activatedOn NetworkValue) (*Version, error) {
	if activatedOn != "" {
		activatedOn = "?activatedOn=" + activatedOn
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/latest%s",
			versions.PropertyID,
			activatedOn,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	newVersions := NewVersions()
	if err := client.BodyJSON(res, newVersions); err != nil {
		return nil, err
	}

	return newVersions.Versions.Items[0], nil
}

// NewVersion creates a new version associated with the Versions collection
func (versions *Versions) NewVersion(createFromVersion *Version, useEtagStrict bool) *Version {
	if createFromVersion == nil {
		var err error
		createFromVersion, err = versions.GetLatestVersion("")
		if err != nil {
			return nil
		}
	}

	version := NewVersion(versions)
	version.CreateFromVersion = createFromVersion.PropertyVersion

	versions.Versions.Items = append(versions.Versions.Items, version)

	if useEtagStrict {
		version.CreateFromVersionEtag = createFromVersion.Etag
	}

	return version
}

// Version represents a Property Version
type Version struct {
	client.Resource
	parent                *Versions
	PropertyVersion       int         `json:"propertyVersion,omitempty"`
	UpdatedByUser         string      `json:"updatedByUser,omitempty"`
	UpdatedDate           time.Time   `json:"updatedDate,omitempty"`
	ProductionStatus      StatusValue `json:"productionStatus,omitempty"`
	StagingStatus         StatusValue `json:"stagingStatus,omitempty"`
	Etag                  string      `json:"etag,omitempty"`
	ProductID             string      `json:"productId,omitempty"`
	Note                  string      `json:"note,omitempty"`
	CreateFromVersion     int         `json:"createFromVersion,omitempty"`
	CreateFromVersionEtag string      `json:"createFromVersionEtag,omitempty"`
	RuleFormat            string      `json:"ruleFormat,omitempty"`
}

// NewVersion creates a new Version
func NewVersion(parent *Versions) *Version {
	version := &Version{parent: parent}
	version.Init()

	return version
}

// GetVersion populates a Version
//
// Api Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaversion
// Endpoint: /papi/v1/properties/{propertyId}/versions/{propertyVersion}{?contractId,groupId}
func (version *Version) GetVersion(property *Property, getVersion int) error {
	if getVersion == 0 {
		getVersion = property.LatestVersion
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d",
			property.PropertyID,
			getVersion,
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

	newVersions := NewVersions()
	if err := client.BodyJSON(res, newVersions); err != nil {
		return err
	}

	version.PropertyVersion = newVersions.Versions.Items[0].PropertyVersion
	version.UpdatedByUser = newVersions.Versions.Items[0].UpdatedByUser
	version.UpdatedDate = newVersions.Versions.Items[0].UpdatedDate
	version.ProductionStatus = newVersions.Versions.Items[0].ProductionStatus
	version.StagingStatus = newVersions.Versions.Items[0].StagingStatus
	version.Etag = newVersions.Versions.Items[0].Etag
	version.ProductID = newVersions.Versions.Items[0].ProductID
	version.Note = newVersions.Versions.Items[0].Note
	version.CreateFromVersion = newVersions.Versions.Items[0].CreateFromVersion
	version.CreateFromVersionEtag = newVersions.Versions.Items[0].CreateFromVersionEtag

	return nil
}

// HasBeenActivated determines if a given version has been activated, optionally on a specific network
func (version *Version) HasBeenActivated(activatedOn NetworkValue) (bool, error) {
	properties := NewProperties()
	property := NewProperty(properties)
	property.PropertyID = version.parent.PropertyID

	property.Group = NewGroup(NewGroups())
	property.Group.GroupID = version.parent.GroupID

	property.Contract = NewContract(NewContracts())
	property.Contract.ContractID = version.parent.ContractID

	activations, err := property.GetActivations()
	if err != nil {
		return false, err
	}

	for _, activation := range activations.Activations.Items {
		if activation.PropertyVersion == version.PropertyVersion && (activatedOn == "" || activation.Network == activatedOn) {
			return true, nil
		}
	}

	return false, nil
}

// Save creates a new version
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#createanewversion
// Endpoint: POST /papi/v1/properties/{propertyId}/versions/{?contractId,groupId}
func (version *Version) Save() error {
	if version.PropertyVersion != 0 {
		return fmt.Errorf("version (%d) already exists", version.PropertyVersion)
	}

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions",
			version.parent.PropertyID,
		),
		version,
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

	var location client.JSONBody
	if err = client.BodyJSON(res, &location); err != nil {
		return err
	}

	req, err = client.NewRequest(
		Config,
		"GET",
		location["versionLink"].(string),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequest(req, true)

	res, err = client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	versions := NewVersions()
	if err = client.BodyJSON(res, versions); err != nil {
		return err
	}

	version.PropertyVersion = versions.Versions.Items[0].PropertyVersion
	version.UpdatedByUser = versions.Versions.Items[0].UpdatedByUser
	version.UpdatedDate = versions.Versions.Items[0].UpdatedDate
	version.ProductionStatus = versions.Versions.Items[0].ProductionStatus
	version.StagingStatus = versions.Versions.Items[0].StagingStatus
	version.Etag = versions.Versions.Items[0].Etag
	version.ProductID = versions.Versions.Items[0].ProductID
	version.Note = versions.Versions.Items[0].Note
	version.CreateFromVersion = versions.Versions.Items[0].CreateFromVersion
	version.CreateFromVersionEtag = versions.Versions.Items[0].CreateFromVersionEtag

	version.parent.AddVersion(version)

	return nil
}
