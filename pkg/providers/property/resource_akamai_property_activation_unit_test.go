package property

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
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
		init              func(*papi.Mock)
		withError         error
	}{
		"version present and set": {
			propertyID:        "prp_id",
			network:           papi.ActivationNetworkStaging,
			versionData:       3,
			versionDataExists: true,
			init:              func(m *papi.Mock) {},
		},
		"version not present but fetched from API": {
			propertyID: "prp_id",
			network:    papi.ActivationNetworkStaging,
			init: func(m *papi.Mock) {
				m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_id",
					ActivatedOn: fmt.Sprintf("%v", papi.ActivationNetworkStaging),
				}).Return(
					&papi.GetPropertyVersionsResponse{
						Version: papi.PropertyVersionGetItem{
							PropertyVersion: 1,
						},
					}, nil).Once()
			},
		},
		"version not present & not fetched - error": {
			propertyID: "prp_id",
			network:    papi.ActivationNetworkProduction,
			init: func(m *papi.Mock) {
				m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_id",
					ActivatedOn: fmt.Sprintf("%v", papi.ActivationNetworkProduction),
				}).Return(nil, tf.ErrNotFound).Once()
			},
			withError: tf.ErrNotFound,
		},
	}
	for name, test := range tests {
		d := schema.TestResourceDataRaw(t, akamaiPropertyActivationSchema, nil)
		if test.versionDataExists {
			d = schema.TestResourceDataRaw(t, akamaiPropertyActivationSchema, map[string]interface{}{
				"version": test.versionData,
			})
		}
		ctx := ctxt{}
		client := &papi.Mock{}
		if test.init != nil {
			test.init(client)
		}
		t.Run(name, func(t *testing.T) {
			version, err := resolveVersion(ctx, d, client, test.propertyID, test.network)
			if test.withError != nil {
				assert.Equal(t, test.withError, err)
				assert.Equal(t, 0, version)
			} else {
				require.NoError(t, err)
				if test.versionDataExists {
					assert.Equal(t, test.versionData, version)
				}
				assert.Less(t, 0, version)
			}
		})
	}
}

func TestLookupActivation(t *testing.T) {

	tests := map[string]struct {
		init                     func(*papi.Mock)
		query                    lookupActivationRequest
		mostRecentActivationDate string
		expectedActivation       *papi.Activation
		expectedError            error
	}{
		"ok": {
			init: func(m *papi.Mock) {
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
			init: func(m *papi.Mock) {
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
			init: func(m *papi.Mock) {
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
			expectedError:            date.ErrDateFormat,
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
			client := &papi.Mock{}
			test.init(client)
			query := lookupActivationRequest{
				propertyID:     "prp_1234",
				version:        2,
				network:        "STAGING",
				activationType: map[papi.ActivationType]struct{}{"ACTIVATE": {}},
			}
			activation, err := lookupActivation(context.TODO(), client, query)
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

func TestCreateActivation(t *testing.T) {

	defer func(t time.Duration) {
		CreateActivationRetry = t
	}(CreateActivationRetry) // restore previous value
	CreateActivationRetry = time.Microsecond

	propID := "someID"

	createReq := papi.CreateActivationRequest{
		PropertyID: propID,
		Activation: papi.Activation{
			ActivationType:  papi.ActivationTypeActivate,
			PropertyID:      "someID",
			PropertyVersion: 1,
			Network:         papi.ActivationNetworkProduction,
		},
	}

	t.Run("create 201 no get", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(&papi.CreateActivationResponse{
			ActivationID: "atv_123",
		}, nil)

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)

	})

	t.Run("create 500 get ok pending", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})

	t.Run("create 400 nonretryable", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 400})).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.NotNil(t, diagErr)
		assert.Equal(t, "", actID)
	})

	t.Run("create 500 retry on empty get", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil)

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})
	t.Run("create 500 retry 3 times before ok", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Times(3)
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 2,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Times(3)

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:27:55Z",
					},
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 2,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)

	})
	t.Run("create 422 get valid", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 422})).Once()

		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})
	t.Run("create 409 get valid", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 409})).Once()

		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})
	t.Run("create 400 return err ", func(t *testing.T) {
		m := &papi.Mock{}

		expectedError := fmt.Errorf("some err: %w", error(&papi.Error{StatusCode: 400}))

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, expectedError).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.NotEmpty(t, diagErr)
		assert.Contains(t, diagErr[0].Summary, expectedError.Error())
		assert.Equal(t, "", actID)
	})
	t.Run("create 500 retry on unexpected version", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 2,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:27:55Z",
					},
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 2,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_000", actID)
	})
	t.Run("create 500 retry on failed status", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()

		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusFailed,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()

		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:27:55Z",
					},
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusFailed,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_000", actID)
	})

	t.Run("create 500 retry on network missmatch", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:27:55Z",
					},
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})
	t.Run("create 500 retry on type missmatch", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:27:55Z",
					},
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})

	t.Run("expect list is sorted and network is filtered", func(t *testing.T) {
		m := &papi.Mock{}

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				// expect sorted -> tested by: mixed elements order
				// expect filter correct network -> tested by: first is incorrect network
				// expect correct type -> tested by: first is incorrect type (should retry)
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T08:43:33Z",
					},
					{
						ActivationID:    "atv_002",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
					{
						ActivationID:    "atv_001",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:20:21Z",
					},
				},
			},
		}, nil).Once()

		m.On("CreateActivation", mock.Anything, createReq).Return(nil, error(&papi.Error{StatusCode: 500})).Once()
		m.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
			PropertyID: propID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:    "atv_000",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T08:43:33Z",
					},
					{
						ActivationID:    "atv_002",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:24:29Z",
					},
					{
						ActivationID:    "atv_001",
						ActivationType:  papi.ActivationTypeDeactivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusActive,
						UpdateDate:      "2023-11-28T13:20:21Z",
					},
					{
						ActivationID:    "atv_123",
						ActivationType:  papi.ActivationTypeActivate,
						PropertyID:      propID,
						PropertyVersion: 1,
						Network:         papi.ActivationNetworkProduction,
						Status:          papi.ActivationStatusPending,
						UpdateDate:      "2023-11-28T13:55:55Z",
					},
				},
			},
		}, nil).Once()

		ctx := context.Background()
		actID, diagErr := createActivation(ctx, m, createReq)
		assert.Nil(t, diagErr)
		assert.Equal(t, "atv_123", actID)
	})
}
