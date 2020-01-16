package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Resources []Resource

type Resource struct {
	APIResourceID           int     `json:"apiResourceId"`
	APIResourceName         string  `json:"apiResourceName"`
	ResourcePath            string  `json:"resourcePath"`
	Description             string  `json:"description"`
	LockVersion             int     `json:"lockVersion"`
	APIResourceClonedFromID *int    `json:"apiResourceClonedFromId"`
	APIResourceLogicID      int     `json:"apiResourceLogicId"`
	CreatedBy               string  `json:"createdBy"`
	CreateDate              string  `json:"createDate"`
	UpdatedBy               string  `json:"updatedBy"`
	UpdateDate              string  `json:"updateDate"`
	APIResourceMethods      Methods `json:"apiResourceMethods"`
}

type ResourceBaseInfo struct {
	APIResourceClonedFromID *int    `json:"apiResourceClonedFromId"`
	APIResourceID           int     `json:"apiResourceId"`
	APIResourceLogicID      int     `json:"apiResourceLogicId"`
	APIResourceName         string  `json:"apiResourceName"`
	CreateDate              string  `json:"createDate"`
	CreatedBy               string  `json:"createdBy"`
	Description             *string `json:"description"`
	Link                    *string `json:"link"`
	LockVersion             int     `json:"lockVersion"`
	Private                 bool    `json:"private"`
	ResourcePath            string  `json:"resourcePath"`
	UpdateDate              string  `json:"updateDate"`
	UpdatedBy               string  `json:"updatedBy"`
}

type ResourceSettings struct {
	Path                 string        `json:"path"`
	Methods              []MethodValue `json:"methods"`
	InheritsFromEndpoint bool          `json:"inheritsFromEndpoint"`
}

func GetResources(endpointId int, version int) (*Resources, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/resources",
			endpointId,
			version,
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

	rep := &Resources{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}
