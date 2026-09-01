package tools

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCPCodeByName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cpCodeName string
		contractID string
		groupID    string
		response   *papi.GetCPCodesResponse
		expected   papi.CPCode
		expectErr  error
	}{
		"single match": {
			cpCodeName: "test cpcode",
			contractID: "ctr_1",
			groupID:    "grp_2",
			response: &papi.GetCPCodesResponse{
				CPCodes: papi.CPCodeItems{Items: []papi.CPCode{
					{ID: "cpc_123", Name: "wrong CP code"},
					{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_1"}},
				}},
			},
			expected: papi.CPCode{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_1"}},
		},
		"no match": {
			cpCodeName: "nonexistent",
			contractID: "ctr_1",
			groupID:    "grp_2",
			response: &papi.GetCPCodesResponse{
				CPCodes: papi.CPCodeItems{Items: []papi.CPCode{
					{ID: "cpc_123", Name: "wrong CP code"},
				}},
			},
			expectErr: ErrCPCodeNotFound,
		},
		"more than one match": {
			cpCodeName: "duplicate",
			contractID: "ctr_1",
			groupID:    "grp_2",
			response: &papi.GetCPCodesResponse{
				CPCodes: papi.CPCodeItems{Items: []papi.CPCode{
					{ID: "cpc_123", Name: "duplicate"},
					{ID: "cpc_234", Name: "duplicate"},
				}},
			},
			expectErr: ErrMoreCPCodesFound,
		},
		"empty list": {
			cpCodeName: "test",
			contractID: "ctr_1",
			groupID:    "grp_2",
			response: &papi.GetCPCodesResponse{
				CPCodes: papi.CPCodeItems{Items: []papi.CPCode{}},
			},
			expectErr: ErrCPCodeNotFound,
		},
		"API error": {
			cpCodeName: "test",
			contractID: "ctr_1",
			groupID:    "grp_2",
			expectErr:  fmt.Errorf("API failure"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &papi.Mock{}

			var apiErr error
			if tc.response == nil {
				apiErr = tc.expectErr
			}
			client.On("GetCPCodes", context.Background(), papi.GetCPCodesRequest{
				ContractID: tc.contractID,
				GroupID:    tc.groupID,
			}).Return(tc.response, apiErr)

			result, err := FindCPCodeByName(context.Background(), client, tc.cpCodeName, tc.contractID, tc.groupID)

			if tc.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			client.AssertExpectations(t)
		})
	}
}
