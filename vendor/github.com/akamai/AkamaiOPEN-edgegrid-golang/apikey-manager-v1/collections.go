package apikeymanager

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Collections []Collection

type Collection struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	KeyCount    int      `json:"keyCount,omitempty"`
	Dirty       bool     `json:"dirty,omitempty"`
	ContractId  string   `json:"contractId,omitempty"`
	GroupId     int      `json:"groupId,omitempty"`
	GrantedACL  []string `json:"grantedACL,omitempty"`
	DirtyACL    []string `json:"dirtyACL,omitempty"`
	Quota       Quota    `json:"quota,omitempty"`
}

func ListCollections() (*Collections, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		"/apikey-manager-api/v1/collections",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collections{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type CreateCollectionOptions struct {
	ContractId  string `json:"contractId,omitempty"`
	GroupId     int    `json:"groupId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func CreateCollection(options *CreateCollectionOptions) (*Collection, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/collections",
		options,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func GetCollection(collectionId int) (*Collection, error) {
	req, err := client.NewJSONRequest(
		Config,
		"GET",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d", collectionId),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func CollectionAclAllow(collectionId int, acl []string) (*Collection, error) {
	collection, err := GetCollection(collectionId)
	if err != nil {
		return collection, err
	}

	acl = append(acl, collection.GrantedACL...)

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d/acl", collectionId),
		acl,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

func CollectionAclDeny(collectionId int, acl []string) (*Collection, error) {
	collection, err := GetCollection(collectionId)
	if err != nil {
		return collection, err
	}

	for cIndex, currentAcl := range collection.GrantedACL {
		for _, newAcl := range acl {
			if newAcl == currentAcl {
				collection.GrantedACL = append(
					collection.GrantedACL[:cIndex],
					collection.GrantedACL[cIndex+1:]...,
				)
			}
		}
	}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d/acl", collectionId),
		collection.GrantedACL,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type Quota struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Value    int    `json:"value,omitempty"`
	Interval string `json:"interval,omitempty"`
	Headers  struct {
		DenyLimitHeaderShown      bool `json:"denyLimitHeaderShown,omitempty"`
		DenyRemainingHeaderShown  bool `json:"denyRemainingHeaderShown,omitempty"`
		DenyNextHeaderShown       bool `json:"denyNextHeaderShown,omitempty"`
		AllowLimitHeaderShown     bool `json:"allowLimitHeaderShown,omitempty"`
		AllowRemainingHeaderShown bool `json:"allowRemainingHeaderShown,omitempty"`
		AllowResetHeaderShown     bool `json:"allowResetHeaderShown,omitempty"`
	} `json:"headers,omitempty"`
}

func CollectionSetQuota(collectionId int, value int) (*Collection, error) {
	collection, err := GetCollection(collectionId)
	if err != nil {
		return collection, err
	}

	collection.Quota.Value = value
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/apikey-manager-api/v1/collections/%d/quota", collectionId),
		collection.Quota,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Collection{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}
