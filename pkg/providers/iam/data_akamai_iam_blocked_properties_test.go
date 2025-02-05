package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestBlockedPropertiesDataSource(t *testing.T) {
	propertyName1 := "example1.com"
	propertyName2 := "example2.com"
	propertyID1 := "prp_123456"
	propertyID2 := "prp_456789"

	tests := map[string]struct {
		givenTF       string
		init          func(*iam.Mock, *papi.Mock)
		expectedCheck resource.TestCheckFunc
		expectError   *regexp.Regexp
	}{
		"happy path - blocked properties are returned": {
			givenTF: "valid.tf",
			init: func(im *iam.Mock, pm *papi.Mock) {
				im.On("ListBlockedProperties", testutils.MockContext, iam.ListBlockedPropertiesRequest{
					IdentityID: "user123",
					GroupID:    1,
				}).Return([]int64{123, 456}, nil)
				im.On("MapPropertyIDToName", testutils.MockContext, iam.MapPropertyIDToNameRequest{
					GroupID:    1,
					PropertyID: 123,
				}).Return(&propertyName1, nil)
				pm.On("MapPropertyNameToID", testutils.MockContext, papi.MapPropertyNameToIDRequest{
					GroupID:    "grp_1",
					ContractID: "ctr_C-123",
					Name:       propertyName1,
				}).Return(&propertyID1, nil)
				im.On("MapPropertyIDToName", testutils.MockContext, iam.MapPropertyIDToNameRequest{
					GroupID:    1,
					PropertyID: 456,
				}).Return(&propertyName2, nil)
				pm.On("MapPropertyNameToID", testutils.MockContext, papi.MapPropertyNameToIDRequest{
					GroupID:    "grp_1",
					ContractID: "ctr_C-123",
					Name:       propertyName2,
				}).Return(&propertyID2, nil)

			},
			expectedCheck: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.#", "2"),
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.0.property_id", "prp_123456"),
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.0.asset_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.1.property_id", "prp_456789"),
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.1.asset_id", "456"),
			),
			expectError: nil,
		},
		"happy path - no blocked properties are returned": {
			givenTF: "valid.tf",
			init: func(im *iam.Mock, _ *papi.Mock) {
				im.On("ListBlockedProperties", testutils.MockContext, iam.ListBlockedPropertiesRequest{
					IdentityID: "user123",
					GroupID:    1,
				}).Return([]int64{}, nil)
			},
			expectedCheck: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", "blocked_properties.#", "0"),
			),
			expectError: nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(im *iam.Mock, _ *papi.Mock) {
				im.On("ListBlockedProperties", testutils.MockContext, iam.ListBlockedPropertiesRequest{
					IdentityID: "user123",
					GroupID:    1,
				}).Return(nil, fmt.Errorf("oops"))
			},
			expectError: regexp.MustCompile("oops"),
		},
		"missing required argument group_id": {
			givenTF:     "missing_group_id.tf",
			expectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
		},
		"missing required argument contract_id": {
			givenTF:     "missing_contract_id.tf",
			expectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
		},
		"missing required argument ui_identity_id": {
			givenTF:     "missing_ui_identity_id.tf",
			expectError: regexp.MustCompile(`The argument "ui_identity_id" is required, but no definition was found`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			papiClient := &papi.Mock{}

			if tc.init != nil {
				tc.init(client, papiClient)
			}
			useIAMandPAPIClient(client, papiClient, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataBlockedProperties/%s", tc.givenTF),
						Check:       tc.expectedCheck,
						ExpectError: tc.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
