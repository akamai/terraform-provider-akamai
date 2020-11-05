package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Search contains SearchProperty method used for fetching properties
	// https://developer.akamai.com/api/core_features/property_manager/v1.html#searchgroup
	Search interface {
		// Search earches properties by name, or by the hostname or edge hostname for which it’s currently active
		// https://developer.akamai.com/api/core_features/property_manager/v1.html#postfindbyvalue
		SearchProperties(context.Context, SearchRequest) (*SearchResponse, error)
	}

	// SearchResponse contains response body of POST /search request
	SearchResponse struct {
		Versions SearchItems `json:"versions"`
	}

	// SearchItems contains a list of search results
	SearchItems struct {
		Items []SearchItem `json:"items"`
	}

	// SearchItem contains details of a search result
	SearchItem struct {
		AccountID        string `json:"accountId"`
		AssetID          string `json:"assetId"`
		ContractID       string `json:"contractId"`
		EdgeHostname     string `json:"edgeHostname"`
		GroupID          string `json:"groupId"`
		Hostname         string `json:"hostname"`
		ProductionStatus string `json:"productionStatus"`
		PropertyID       string `json:"propertyId"`
		PropertyName     string `json:"propertyName"`
		PropertyVersion  int    `json:"propertyVersion"`
		StagingStatus    string `json:"stagingStatus"`
		UpdatedByUser    string `json:"updatedByUser"`
		UpdatedDate      string `json:"updatedDate"`
	}

	// SearchRequest contains key-value pair for search request
	// Key must have one of three values: "edgeHostname", "hostname" or "propertyName"
	SearchRequest struct {
		Key   string
		Value string
	}
)

const (
	// SearchKeyEdgeHostname search request key
	SearchKeyEdgeHostname = "edgeHostname"
	// SearchKeyHostname search request key
	SearchKeyHostname = "hostname"
	// SearchKeyPropertyName search request key
	SearchKeyPropertyName = "propertyName"
)

// Validate validate SearchRequest struct
func (s SearchRequest) Validate() error {
	return validation.Errors{
		"SearchKey": validation.Validate(s.Key,
			validation.Required,
			validation.In(SearchKeyEdgeHostname, SearchKeyHostname, SearchKeyPropertyName)),
		"SearchValue": validation.Validate(s.Value, validation.Required),
	}.Filter()
}

var (
	ErrSearchProperties = errors.New("searching for properties")
)

func (p *papi) SearchProperties(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrSearchProperties, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("SearchProperties")

	searchURL := "/papi/v1/search/find-by-value"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrSearchProperties, err)
	}

	var search SearchResponse
	resp, err := p.Exec(req, &search, map[string]string{request.Key: request.Value})
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrSearchProperties, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrSearchProperties, p.Error(resp))
	}

	return &search, nil
}
