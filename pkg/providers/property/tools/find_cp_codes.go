// Package tools for common property provider operations
package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
)

// ErrMoreCPCodesFound is returned when more than one CP code was found for the given name
var ErrMoreCPCodesFound = errors.New("more than one CP code was found for the given name")

// ErrCPCodeNotFound is returned when no CP code was found for the given name
var ErrCPCodeNotFound = errors.New("CP code not found")

// FindCPCodeByName searches all CP codes for a match against given name
func FindCPCodeByName(ctx context.Context, client papi.PAPI, name, contractID, groupID string) (papi.CPCode, error) {
	r, err := client.GetCPCodes(ctx, papi.GetCPCodesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return papi.CPCode{}, err
	}

	var matchedCPCodes []papi.CPCode
	for _, cpc := range r.CPCodes.Items {
		if cpc.Name == name {
			matchedCPCodes = append(matchedCPCodes, cpc)
		}
	}

	if len(matchedCPCodes) > 1 {
		return papi.CPCode{}, fmt.Errorf("%w: more than one CP code for name %s was found", ErrMoreCPCodesFound, name)
	} else if len(matchedCPCodes) == 1 {
		return matchedCPCodes[0], nil
	}

	return papi.CPCode{}, fmt.Errorf("%w: CP code: %s", ErrCPCodeNotFound, name)
}
