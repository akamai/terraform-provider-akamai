package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	// Contracts contains operations available on Contract resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#contractsgroup
	Contracts interface {
		// GetContract provides a read-only list of contract names and identifiers
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getcontracts
		GetContracts(context.Context) (*GetContractsResponse, error)
	}

	// Contract represents a property contract resource
	Contract struct {
		ContractID       string `json:"contractId"`
		ContractTypeName string `json:"contractTypeName"`
	}

	// ContractsItems is the response items array
	ContractsItems struct {
		Items []*Contract `json:"items"`
	}

	// GetContractsResponse represents a collection of property manager contracts
	// This is the reponse to the /papi/v1/contracts request
	GetContractsResponse struct {
		AccountID string         `json:"accountId"`
		Contracts ContractsItems `json:"contracts"`
	}
)

var (
	ErrGetContracts = errors.New("fetching contracts")
)

func (p *papi) GetContracts(ctx context.Context) (*GetContractsResponse, error) {
	var contracts GetContractsResponse

	logger := p.Log(ctx)
	logger.Debug("GetContracts")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/papi/v1/contracts", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetContracts, err)
	}

	resp, err := p.Exec(req, &contracts)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetContracts, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetContracts, p.Error(resp))
	}

	return &contracts, nil
}
