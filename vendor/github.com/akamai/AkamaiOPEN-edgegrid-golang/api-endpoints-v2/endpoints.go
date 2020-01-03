package apiendpoints

import (
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/google/go-querystring/query"
)

type Endpoints []Endpoint

type Endpoint struct {
	APICategoryIds             []int                 `json:"apiCategoryIds,omitempty"`
	APIEndPointHosts           []string              `json:"apiEndPointHosts"`
	APIEndPointID              int                   `json:"apiEndPointId,omitempty"`
	APIEndPointLocked          bool                  `json:"apiEndPointLocked,omitempty"`
	APIEndPointName            string                `json:"apiEndPointName"`
	APIEndPointScheme          string                `json:"apiEndPointScheme,omitempty"`
	APIResourceBaseInfo        []*ResourceBaseInfo   `json:"apiResourceBaseInfo,omitempty"`
	BasePath                   string                `json:"basePath,omitempty"`
	ClonedFromVersion          *int                  `json:"clonedFromVersion,omitempty"`
	ConsumeType                string                `json:"consumeType,omitempty"`
	ContractID                 string                `json:"contractId,omitempty"`
	CreateDate                 string                `json:"createDate,omitempty"`
	CreatedBy                  string                `json:"createdBy,omitempty"`
	Description                string                `json:"description,omitempty"`
	GroupID                    int                   `json:"groupId,omitempty"`
	ProductionVersion          *VersionSummary       `json:"productionVersion,omitempty"`
	ProductionStatus           string                `json:"productionStatus,omitempty"`
	ProtectedByAPIKey          bool                  `json:"protectedByApiKey,omitempty"`
	StagingStatus              string                `json:"stagingStatus,omitempty"`
	StagingVersion             *VersionSummary       `json:"stagingVersion,omitempty"`
	UpdateDate                 string                `json:"updateDate,omitempty"`
	UpdatedBy                  string                `json:"updatedBy,omitempty"`
	VersionNumber              int                   `json:"versionNumber,omitempty"`
	SecurityScheme             *SecurityScheme       `json:"securityScheme,omitempty"`
	AkamaiSecurityRestrictions *SecurityRestrictions `json:"akamaiSecurityRestrictions,omitempty"`
	APIResources               *Resources            `json:"apiResources,omitempty"`
}

type SecurityScheme struct {
	SecuritySchemeType   string                `json:"securitySchemeType,omitempty"`
	SecuritySchemeDetail *SecuritySchemeDetail `json:"securitySchemeDetail,omitempty"`
}

type SecuritySchemeDetail struct {
	APIKeyLocation string `json:"apiKeyLocation,omitempty"`
	APIKeyName     string `json:"apiKeyName,omitempty"`
}

type SecurityRestrictions struct {
	MaxJsonxmlElement       int `json:"MAX_JSONXML_ELEMENT,omitempty"`
	MaxElementNameLength    int `json:"MAX_ELEMENT_NAME_LENGTH,omitempty"`
	MaxDocDepth             int `json:"MAX_DOC_DEPTH,omitempty"`
	PositiveSecurityEnabled int `json:"POSITIVE_SECURITY_ENABLED,omitempty"`
	MaxStringLength         int `json:"MAX_STRING_LENGTH,omitempty"`
	MaxBodySize             int `json:"MAX_BODY_SIZE,omitempty"`
	MaxIntegerValue         int `json:"MAX_INTEGER_VALUE,omitempty"`
}

type CreateEndpointOptions struct {
	ContractId string   `json:"contractId,omitempty"`
	GroupId    int      `json:"groupId,omitempty"`
	Name       string   `json:"apiEndPointName,omitempty"`
	BasePath   string   `json:"basePath,omitempty"`
	Hostnames  []string `json:"apiEndPointHosts,omitempty"`
}

func CreateEndpoint(options *CreateEndpointOptions) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/api-definitions/v2/endpoints",
		options,
	)

	return call(req, err)
}

type CreateEndpointFromFileOptions struct {
	File       string
	Format     string
	ContractId string
	GroupId    int
}

func CreateEndpointFromFile(options *CreateEndpointFromFileOptions) (*Endpoint, error) {
	req, err := client.NewMultiPartFormDataRequest(
		Config,
		"/api-definitions/v2/endpoints/files",
		options.File,
		map[string]string{
			"contractId":       options.ContractId,
			"groupId":          strconv.Itoa(options.GroupId),
			"importFileFormat": options.Format,
		},
	)

	return call(req, err)
}

type UpdateEndpointFromFileOptions struct {
	EndpointId int
	Version    int
	File       string
	Format     string
}

func UpdateEndpointFromFile(options *UpdateEndpointFromFileOptions) (*Endpoint, error) {
	url := fmt.Sprintf(
		"/api-definitions/v2/endpoints/%d/versions/%d/file",
		options.EndpointId,
		options.Version,
	)

	req, err := client.NewMultiPartFormDataRequest(
		Config,
		url,
		options.File,
		map[string]string{
			"importFileFormat": options.Format,
		},
	)

	return call(req, err)
}

type ListEndpointOptions struct {
	ContractId        string `url:"contractId,omitempty"`
	GroupId           int    `url:"groupId,omitempty"`
	Category          string `url:"category,omitempty"`
	Contains          string `url:"contains,omitempty"`
	Page              int    `url:"page,omitempty"`
	PageSize          int    `url:"pageSize,omitempty"`
	Show              string `url:show,omitempty`
	SortBy            string `url:"sortBy,omitempty"`
	SortOrder         string `url:"sortOrder,omitempty"`
	VersionPreference string `url:"versionPreference,omitempty"`
}

type EndpointList struct {
	APIEndPoints Endpoints `json:"apiEndPoints"`
	Links        Links     `json:"links"`
	Page         int       `json:"page"`
	PageSize     int       `json:"pageSize"`
	TotalSize    int       `json:"totalSize"`
}

func (list *EndpointList) ListEndpoints(options *ListEndpointOptions) error {
	q, err := query.Values(options)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"/api-definitions/v2/endpoints?%s",
		q.Encode(),
	)

	req, err := client.NewJSONRequest(Config, "GET", url, nil)
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

	if err = client.BodyJSON(res, list); err != nil {
		return err
	}

	return nil
}

func RemoveEndpoint(endpointId int) (*Endpoint, error) {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d",
			endpointId,
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

	rep := &Endpoint{}
	return rep, nil
}
