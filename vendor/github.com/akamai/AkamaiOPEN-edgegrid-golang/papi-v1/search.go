package papi

import (
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

type SearchKey string

const (
	SearchByPropertyName SearchKey = "propertyName"
	SearchByHostname     SearchKey = "hostname"
	SearchByEdgeHostname SearchKey = "edgeHostname"
)

type SearchResult struct {
	Versions struct {
		Items []struct {
			UpdatedByUser    string    `json:"updatedByUser"`
			StagingStatus    string    `json:"stagingStatus"`
			AssetID          string    `json:"assetId"`
			PropertyName     string    `json:"propertyName"`
			PropertyVersion  int       `json:"propertyVersion"`
			UpdatedDate      time.Time `json:"updatedDate"`
			ContractID       string    `json:"contractId"`
			AccountID        string    `json:"accountId"`
			GroupID          string    `json:"groupId"`
			PropertyID       string    `json:"propertyId"`
			ProductionStatus string    `json:"productionStatus"`
		} `json:"items"`
	} `json:"versions"`
}

// Search searches for properties
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#postfindbyvalue
// Endpoint: POST /papi/v1/search/find-by-value
func Search(searchBy SearchKey, propertyName string, correlationid string) (*SearchResult, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/papi/v1/search/find-by-value",
		map[string]string{(string)(searchBy): propertyName},
	)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	results := &SearchResult{}
	if err = client.BodyJSON(res, results); err != nil {
		return nil, err
	}

	if len(results.Versions.Items) == 0 {
		return nil, nil
	}

	return results, nil
}
