package property

import (
	"errors"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"testing"
)

func TestLookupActivation(t *testing.T) {

	tests := map[string]struct {
		init                     func(*mockpapi)
		query                    lookupActivationRequest
		mostRecentActivationDate string
		expectedActivation       *papi.Activation
		expectedError            error
	}{
		"ok": {
			init: func(m *mockpapi) {
				m.On("GetActivations", mock.Anything, mock.Anything).Return(
					&papi.GetActivationsResponse{
						Response: papi.Response{
							AccountID:  "act_1234",
							ContractID: "ctr_1234",
							GroupID:    "grp_1234",
							Etag:       "1234",
							Errors:     nil,
							Warnings:   nil,
						},
						Activations: papi.ActivationsItems{Items: []*papi.Activation{
							{
								Status: "ABORTED",
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 1,
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
								Network:         "PRODUCTION",
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
								Network:         "STAGING",
								SubmitDate:      "2006-01-02T15:04:05Z",
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
								Network:         "STAGING",
								SubmitDate:      "2016-03-22T15:04:05Z",
								ActivationID:    "1234",
							},
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
								Network:         "STAGING",
								SubmitDate:      "2014-03-22T15:04:05Z",
							},
						}},
					}, nil,
				)
			},
			query: lookupActivationRequest{
				propertyID:     "prp_1234",
				version:        2,
				network:        "STAGING",
				activationType: map[papi.ActivationType]struct{}{"ACTIVATE": {}},
			},
			mostRecentActivationDate: "2016-03-22T15:04:05Z",
			expectedError:            nil,
			expectedActivation: &papi.Activation{
				AccountID:       "act_1234",
				ActivationID:    "1234",
				ActivationType:  "ACTIVATE",
				GroupID:         "grp_1234",
				PropertyID:      "prp_1234",
				PropertyVersion: 2,
				Network:         "STAGING",
				Status:          "ACTIVE",
				SubmitDate:      "2016-03-22T15:04:05Z",
			},
		},
		"ok, but no activations": {
			init: func(m *mockpapi) {
				m.On("GetActivations", mock.Anything, mock.Anything).Return(
					&papi.GetActivationsResponse{
						Response: papi.Response{
							AccountID:  "act_1234",
							ContractID: "ctr_1234",
							GroupID:    "grp_1234",
							Etag:       "1234",
							Errors:     nil,
							Warnings:   nil,
						},
						Activations: papi.ActivationsItems{Items: []*papi.Activation{
							{
								Status: "ABORTED",
							},
						}},
					}, nil,
				)
			},
			query: lookupActivationRequest{
				propertyID:     "prp_1234",
				version:        2,
				network:        "STAGING",
				activationType: map[papi.ActivationType]struct{}{"ACTIVATE": {}},
			},
			mostRecentActivationDate: "2016-03-22T15:04:05Z",
			expectedError:            nil,
			expectedActivation:       nil,
		},
		"date parse error": {
			init: func(m *mockpapi) {
				m.On("GetActivations", mock.Anything, mock.Anything).Return(
					&papi.GetActivationsResponse{
						Response: papi.Response{
							AccountID:  "act_1234",
							ContractID: "ctr_1234",
							GroupID:    "grp_1234",
							Etag:       "1234",
							Errors:     nil,
							Warnings:   nil,
						},
						Activations: papi.ActivationsItems{Items: []*papi.Activation{
							{
								Status:          "ACTIVE",
								PropertyVersion: 2,
								ActivationType:  "ACTIVATE",
								Network:         "STAGING",
								SubmitDate:      "2016-13-22T15:04:05Z",
								ActivationID:    "1234",
							},
						}},
					}, nil,
				)
			},
			query: lookupActivationRequest{
				propertyID:     "prp_1234",
				version:        2,
				network:        "STAGING",
				activationType: map[papi.ActivationType]struct{}{"ACTIVATE": {}},
			},
			mostRecentActivationDate: "2016-03-22T15:04:05Z",
			expectedError:            tools.ErrDateFormat,
			expectedActivation: &papi.Activation{
				AccountID:       "act_1234",
				ActivationID:    "1234",
				ActivationType:  "ACTIVATE",
				GroupID:         "grp_1234",
				PropertyID:      "prp_1234",
				PropertyVersion: 2,
				Network:         "STAGING",
				Status:          "ACTIVE",
				SubmitDate:      "2016-03-22T15:04:05Z",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			query := lookupActivationRequest{
				propertyID:     "prp_1234",
				version:        2,
				network:        "STAGING",
				activationType: map[papi.ActivationType]struct{}{"ACTIVATE": {}},
			}
			activation, err := lookupActivation(nil, client, query)
			assert.True(t, errors.Is(err, test.expectedError))
			if err == nil {
				expectedActivation := test.expectedActivation
				if expectedActivation != nil {
					assert.Equal(t, expectedActivation.ActivationType, activation.ActivationType)
					assert.Equal(t, test.mostRecentActivationDate, activation.SubmitDate)
				} else {
					assert.Nil(t, activation)
				}
			}
		})
	}
}
