package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourcePAPIPropertyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"property activation - OK": {
			init: func(m *papi.Mock) {
				activationsResponseDeactivated = papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
						Note:            "property activation note for creating",
						NotifyEmails:    []string{"some@example.com"},
					}}},
				}
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/ok/datasource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "contact.0", "some@example.com"),
					),
				},
			},
		},
		"check schema property activation - OK": {
			init: func(m *papi.Mock) {
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T14:04:05Z", nil), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/ok/datasource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "note", ""),
					),
				},
			},
		},
		"version not provided: default to latest - OK": {
			init: func(m *papi.Mock) {
				m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_test",
					ActivatedOn: fmt.Sprintf("%v", papi.ActivationNetworkStaging),
				}).Return(&papi.GetPropertyVersionsResponse{
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T14:04:05Z", nil), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/ok/datasource_property_activation_no_version.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_property_activation.test", "note", ""),
					),
				},
			},
		},
		"check schema property activation - papi error": {
			init: func(m *papi.Mock) {
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, fmt.Errorf("failed to create request")).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/ok/datasource_property_activation.tf"),
					ExpectError: regexp.MustCompile("failed to create request"),
				},
			},
		},
		"check schema property activation - no property id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/no_propertyId/datasource_property_activation.tf"),
					ExpectError: regexp.MustCompile("Error: Missing required argument"),
				},
			},
		},
		"no active property on production": {
			init: func(m *papi.Mock) {
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestDSPropertyActivation/ok/datasource_property_activation_prod.tf"),
					ExpectError: regexp.MustCompile("there is no active version on PRODUCTION network"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
