package apikeymanager

import (
	"encoding/json"
	"io/ioutil"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Keys []int

type Key struct {
	Id                  int      `json:"id,omitempty"`
	Value               string   `json:"value,omitempty"`
	Label               string   `json:"label,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	CollectionName      string   `json:"collectionName,omitempty"`
	CollectionId        int      `json:"collectionId,omitempty"`
	Description         string   `json:"description,omitempty"`
	Revoked             bool     `json:"revoked,omitempty"`
	Dirty               bool     `json:"dirty,omitempty"`
	CreatedAt           string   `json:"createdAt,omitempty"`
	RevokedAt           string   `json:"revokedAt,omitempty"`
	TerminationAt       string   `json:"terminationAt,omitempty"`
	QuotaUsage          int      `json:"quotaUsage,omitempty"`
	QuotaUsageTimestamp string   `json:"quotaUsageTimestamp,omitempty"`
	QuotaUpdateState    string   `json:"quotaUpdateState,omitempty"`
}

type CreateKey struct {
	Value        string   `json:"value,omitempty"`
	Label        string   `json:"label,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CollectionId int      `json:"collectionId,omitempty"`
	Description  string   `json:"description,omitempty"`
	Mode         string   `json:"mode,omitempty"`
}

func CollectionAddKey(collectionId int, name, value string) (*Key, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/keys",
		&CreateKey{
			Label:        name,
			Value:        value,
			CollectionId: collectionId,
			Mode:         "CREATE_ONE",
		},
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

	rep := &Key{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type ImportKey struct {
	Name         string `json:"name,omitempty"`
	Content      string `json:"content,omitempty"`
	CollectionId int    `json:"collectionId,omitempty"`
}

func CollectionImportKeys(collectionId int, filename string) (*Keys, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/keys/import",
		&ImportKey{
			Name:         filename,
			CollectionId: collectionId,
			Content:      string(fileContent),
		},
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

	rep := &Keys{}
	err = json.Unmarshal(fileContent, rep)

	return rep, err
}

type RevokeKeys struct {
	Keys Keys `json:"keys,omitempty"`
}

func RevokeKey(key int) (*Key, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/apikey-manager-api/v1/keys/revoke",
		&RevokeKeys{
			Keys: Keys{key},
		},
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

	return &Key{}, nil
}
