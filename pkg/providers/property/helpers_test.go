package property

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/stretchr/testify/assert"
)

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
		"ok prod": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkProduction,
			addNetwork:  "PROD",
		},
		"ok stag": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkStaging,
			addNetwork:  "STAG",
		},
		"ok stage": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkStaging,
			addNetwork:  "STAGE",
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
		t.Run(name, func(t *testing.T) {
			net, err := NetworkAlias(string(test.addNetwork))
			resultNetwork := papi.ActivationNetwork(net)

			if test.withError != nil {
				assert.Error(t, test.withError, err)
			} else {
				assert.Equal(t, test.networkTest, resultNetwork)
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsPropertyInGroup(t *testing.T) {
	key := papiKey{
		propertyID: "prp_1",
		groupID:    "grp_2",
		contractID: "ctr_3",
	}
	tests := map[string]struct {
		init          func(*testing.T, *mockProperty)
		expected      bool
		expectedError string
	}{
		"returns true if property found and group matches": {
			init: func(t *testing.T, p *mockProperty) {
				p.mockGetProperty()
			},
			expected: true,
		},
		"returns false if property found but group differs": {
			init: func(t *testing.T, p *mockProperty) {
				req := p.getPropertyRequest()
				res := p.getPropertyResponse()
				res.Property.GroupID = "grp_555"
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(&res, nil)
			},
			expected: false,
		},
		"returns false if property not found (HTTP 403)": {
			init: func(t *testing.T, p *mockProperty) {
				req := p.getPropertyRequest()
				err := papi.Error{
					StatusCode: http.StatusForbidden,
				}
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(nil, &err)
			},
			expected: false,
		},
		"forwards fetching property error": {
			init: func(t *testing.T, p *mockProperty) {
				req := p.getPropertyRequest()
				err := errors.New("dummy error")
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(nil, err)
			},
			expectedError: "unexpected http error for {prp_1 grp_2 ctr_3}: dummy error",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mock := &papi.Mock{}
			mp := mockProperty{
				mockPropertyData: mockPropertyData{
					propertyID: key.propertyID,
					groupID:    key.groupID,
					contractID: key.contractID,
				},
				papiMock: mock,
			}
			test.init(t, &mp)
			hlp := helper{mock, nil}
			res, err := hlp.isPropertyInGroup(context.Background(), key)
			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}
}

func TestValidatePropertyMove(t *testing.T) {
	key := papiKey{
		propertyID: "prp_1",
		groupID:    "grp_2",
		contractID: "ctr_3",
	}
	tests := map[string]struct {
		init            func(*testing.T, *mockProperty)
		validationError *string
	}{
		"no error for past activations": {
			init: func(t *testing.T, p *mockProperty) {
				p.activations = papi.ActivationsItems{
					Items: []*papi.Activation{
						{
							ActivationID: "atv_1234567",
						},
					},
				}
				p.mockGetActivationsCompleteRequest()
			},
			validationError: nil,
		},
		"validation error for no activations": {
			init: func(t *testing.T, p *mockProperty) {
				p.mockGetActivationsCompleteRequest()
			},
			validationError: ptr.To("moving properties that have never been activated is not supported " +
				"(property id: prp_1, contract id: ctr_3, group id grp_2)"),
		},
		"forwards fetching activations error": {
			init: func(t *testing.T, p *mockProperty) {
				p.mockGetActivationsCompleteRequest(errors.New("dummy error"))
			},
			validationError: ptr.To("error getting activations list for {prp_1 grp_2 ctr_3}: dummy error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mock := &papi.Mock{}
			mp := mockProperty{
				mockPropertyData: mockPropertyData{
					propertyID: key.propertyID,
					groupID:    key.groupID,
					contractID: key.contractID,
				},
				papiMock: mock,
			}
			test.init(t, &mp)
			hlp := helper{mock, nil}
			err := hlp.validatePropertyMove(context.Background(), key)
			if test.validationError != nil {
				assert.ErrorContains(t, err, *test.validationError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWaitForPropertyGroupIDChange(t *testing.T) {
	key := papiKey{
		propertyID: "prp_1",
		groupID:    "grp_2",
		contractID: "ctr_3",
	}
	tests := map[string]struct {
		init          func(*testing.T, *mockProperty)
		expectedError *string
	}{
		"returns immediately if property in group": {
			init: func(t *testing.T, p *mockProperty) {
				p.mockGetProperty()
			},
			expectedError: nil,
		},
		"retries twice until property in group": {
			init: func(t *testing.T, p *mockProperty) {
				// 2x different group
				req := p.getPropertyRequest()
				res := p.getPropertyResponse()
				res.Property.GroupID = "grp_555"
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(&res, nil).Twice()

				// desired group
				p.mockGetProperty()
			},
			expectedError: nil,
		},
		"returns error if all three attempts exhausted": {
			init: func(t *testing.T, p *mockProperty) {
				// 3x different group
				req := p.getPropertyRequest()
				res := p.getPropertyResponse()
				res.Property.GroupID = "grp_555"
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(&res, nil).Times(3)
			},
			expectedError: ptr.To("waiting for groupID change to: grp_2 for propertyID: prp_1, " +
				"contractID: ctr_3 in 3 attempts failed"),
		},
		"forwards fetching property error": {
			init: func(t *testing.T, p *mockProperty) {
				req := p.getPropertyRequest()
				err := errors.New("dummy error")
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(nil, err)
			},
			expectedError: ptr.To("unexpected http error for {prp_1 grp_2 ctr_3}: dummy error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mock := &papi.Mock{}
			mp := mockProperty{
				mockPropertyData: mockPropertyData{
					propertyID: key.propertyID,
					groupID:    key.groupID,
					contractID: key.contractID,
				},
				papiMock: mock,
			}
			test.init(t, &mp)
			hlp := helper{mock, nil}
			err := hlp.waitForPropertyGroupIDChange(context.Background(), key, 3, time.Second)
			if test.expectedError != nil {
				assert.ErrorContains(t, err, *test.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// For general tests of the moving property functionality, see group id scenarios
// in TestPropertyLifecycle.
func TestMovePropertyValidations(t *testing.T) {
	key := func(groupID string) papiKey {
		return papiKey{
			propertyID: "prp_1",
			groupID:    groupID,
			contractID: "ctr_3",
		}
	}
	tests := map[string]struct {
		key           papiKey
		assetID       string
		destGroupID   string
		expectedError string
	}{
		"invalid src group id": {
			key:           key("group_12345"),
			assetID:       "aid_12345678",
			destGroupID:   "grp_67890",
			expectedError: "error parsing src group id: strconv.Atoi: parsing \"group_12345\": invalid syntax",
		},
		"invalid asset id": {
			key:           key("grp_12345"),
			assetID:       "0xa2345678",
			destGroupID:   "grp_67890",
			expectedError: "error parsing asset id: strconv.Atoi: parsing \"0xa2345678\": invalid syntax",
		},
		"invalid dest group id": {
			key:           key("12345"),
			assetID:       "12345678",
			destGroupID:   "67890.4",
			expectedError: "error parsing dst group id: strconv.Atoi: parsing \"67890.4\": invalid syntax",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hlp := helper{}
			err := hlp.moveProperty(context.Background(), test.key, test.assetID, test.destGroupID)
			assert.ErrorContains(t, err, test.expectedError)
		})
	}
}
