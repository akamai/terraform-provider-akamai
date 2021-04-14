package property

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	ctxt struct {
	}
)

func (c ctxt) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (c ctxt) Done() <-chan struct{} {
	panic("implement me")
}

func (c ctxt) Err() error {
	panic("implement me")
}

func (c ctxt) Value(_ interface{}) interface{} {
	panic("implement me")
}

func TestResolveVersion(t *testing.T) {
	tests := map[string]struct {
		versionData       int
		versionDataExists bool
		propertyID        string
		network           papi.ActivationNetwork
		init              func(*mockpapi)
		withError         error
	}{
		"ok": {
			versionData:       0,
			versionDataExists: true,
			init: func(m *mockpapi) {
				m.On("GetLatestVersion", mock.Anything, mock.Anything).Return(
					&papi.GetPropertyVersionsResponse{
						Version: papi.PropertyVersionGetItem{
							PropertyVersion: 0,
						},
					}, nil)
			},
		},
		"version not present but fetched": {
			versionData: 1,
			init: func(m *mockpapi) {
				m.On("GetLatestVersion", mock.Anything, mock.Anything).Return(
					&papi.GetPropertyVersionsResponse{
						Version: papi.PropertyVersionGetItem{
							PropertyVersion: 1,
						},
					}, nil)
			},
		},
		"version not present & not fetched": {
			versionData: 0,
			init: func(m *mockpapi) {
				m.On("GetLatestVersion", mock.Anything, mock.Anything).Return(
					&papi.GetPropertyVersionsResponse{
						Version: papi.PropertyVersionGetItem{
							PropertyVersion: 1,
						},
					}, tools.ErrNotFound)
			},
			withError: tools.ErrNotFound,
		},
	}
	for name, test := range tests {
		d := schema.ResourceData{}
		ctx := ctxt{}
		client := &mockpapi{}
		test.init(client)
		t.Run(name, func(t *testing.T) {
			if test.versionDataExists {
				d.Set("version", test.versionData)
			}

			version, err := resolveVersion(ctx, &d, client, test.propertyID, test.network)
			if test.withError != nil {
				assert.Equal(t, test.withError, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.versionData, version)
		})
	}
}

func TestNetworkAlias(t *testing.T) {
	tests := map[string]struct {
		hasNetwork  bool
		addNetwork  papi.ActivationNetwork
		networkTest papi.ActivationNetwork
		withError   error
	}{
		"ok production": {
			hasNetwork:  true,
			addNetwork:  papi.ActivationNetworkProduction,
			networkTest: papi.ActivationNetworkProduction,
		},
		"ok p": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkProduction,
			addNetwork:  "P",
		},
		"ok default staging": {
			hasNetwork:  false,
			addNetwork:  papi.ActivationNetworkStaging,
			networkTest: papi.ActivationNetworkStaging,
		},
		"nok malformed input": {
			hasNetwork: true,
			addNetwork: "other",
			withError:  fmt.Errorf("network not recognized"),
		},
	}
	for name, test := range tests {
		resource := schema.Resource{
			Schema: map[string]*schema.Schema{
				"network": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  papi.ActivationNetworkStaging,
				},
			},
		}
		d := resource.TestResourceData()
		if test.hasNetwork {
			_ = d.Set("network", string(test.addNetwork))
		}
		t.Run(name, func(t *testing.T) {
			net, err := networkAlias(d)

			if test.withError != nil {
				assert.Error(t, test.withError, err)
			} else {
				assert.NotNil(t, net)
				assert.Equal(t, test.networkTest, net)
			}
		})
	}
}

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
				if test.expectedActivation != nil {
					assert.Equal(t, test.expectedActivation.ActivationType, activation.ActivationType)
					assert.Equal(t, test.mostRecentActivationDate, activation.SubmitDate)
				} else {
					assert.Nil(t, activation)
				}
			}
		})
	}
}
