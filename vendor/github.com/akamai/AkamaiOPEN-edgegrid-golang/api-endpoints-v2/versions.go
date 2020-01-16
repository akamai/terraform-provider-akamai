package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Versions struct {
	APIEndPointID   int       `json:"apiEndPointId"`
	APIEndPointName string    `json:"apiEndPointName"`
	APIVersions     []Version `json:"apiVersions"`
}

type Version struct {
	CreatedBy            string       `json:"createdBy"`
	CreateDate           string       `json:"createDate"`
	UpdateDate           string       `json:"updateDate"`
	UpdatedBy            string       `json:"updatedBy"`
	APIEndPointVersionID int          `json:"apiEndPointVersionId"`
	BasePath             string       `json:"basePath"`
	Description          *string      `json:"description`
	BasedOn              *int         `json:"basedOn"`
	StagingStatus        *StatusValue `json:"stagingStatus"`
	ProductionStatus     *StatusValue `json:"productionStatus"`
	StagingDate          *string      `json:"stagingDate"`
	ProductionDate       *string      `json:"productionDate"`
	IsVersionLocked      bool         `json:"isVersionLocked"`
	AvailableActions     []string     `json:"availableActions"`
	VersionNumber        int          `json:"versionNumber"`
	LockVersion          int          `json:"lockVersion"`
}

type VersionSummary struct {
	Status        StatusValue `json:"status,omitempty"`
	VersionNumber int         `json:"versionNumber,omitempty"`
}

type StatusValue string

const (
	StatusPending     string = "PENDING"
	StatusActive      string = "ACTIVE"
	StatusDeactivated string = "DEACTIVATED"
	StatusFailed      string = "FAILED"
)

type ListVersionsOptions struct {
	EndpointId int
}

func ListVersions(options *ListVersionsOptions) (*Versions, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions",
			options.EndpointId,
		),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Versions{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type GetVersionOptions struct {
	EndpointId int
	Version    int
}

func GetVersion(options *GetVersionOptions) (*Endpoint, error) {
	if options.Version == 0 {
		versions, err := ListVersions(&ListVersionsOptions{EndpointId: options.EndpointId})
		if err != nil {
			return nil, err
		}

		loc := len(versions.APIVersions) - 1
		v := versions.APIVersions[loc]
		options.Version = v.VersionNumber
	}

	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/resources-detail",
			options.EndpointId,
			options.Version,
		),
		nil,
	)

	return call(req, err)
}

func ModifyVersion(endpoint *Endpoint) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d",
			endpoint.APIEndPointID,
			endpoint.VersionNumber,
		),
		endpoint,
	)

	return call(req, err)
}

type CloneVersionOptions struct {
	EndpointId int
	Version    int
}

func CloneVersion(options *CloneVersionOptions) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/cloneVersion",
			options.EndpointId,
			options.Version,
		),
		options,
	)

	return call(req, err)
}

type RemoveVersionOptions struct {
	EndpointId    int
	VersionNumber int
}

func RemoveVersion(options *RemoveVersionOptions) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d",
			options.EndpointId,
			options.VersionNumber,
		),
		nil,
	)

	return call(req, err)
}
