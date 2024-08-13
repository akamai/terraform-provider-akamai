package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestBlockedPropertiesDataSource(t *testing.T) {
	propertyName1 := "example1.com"
	propertyName2 := "example2.com"
	propertyID1 := "prp_123456"
	propertyID2 := "prp_456789"
	tests := map[string]struct {
		givenTF                   string
		init                      func(*iam.Mock, *papi.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"happy path - blocked properties are returned": {
			givenTF: "valid.tf",
			init: func(im *iam.Mock, pm *papi.Mock) {
				im.On("ListBlockedProperties", mock.Anything, iam.ListBlockedPropertiesRequest{
					IdentityID: "user123",
					GroupID:    1,
				}).Return([]int64{123, 456}, nil)
				im.On("MapPropertyIDToName", mock.Anything, iam.MapPropertyIDToNameRequest{
					GroupID:    1,
					PropertyID: 123,
				}).Return(&propertyName1, nil)
				pm.On("MapPropertyNameToID", mock.Anything, papi.MapPropertyNameToIDRequest{
					GroupID:    "grp_1",
					ContractID: "ctr_C-123",
					Name:       propertyName1,
				}).Return(&propertyID1, nil)
				im.On("MapPropertyIDToName", mock.Anything, iam.MapPropertyIDToNameRequest{
					GroupID:    1,
					PropertyID: 456,
				}).Return(&propertyName2, nil)
				pm.On("MapPropertyNameToID", mock.Anything, papi.MapPropertyNameToIDRequest{
					GroupID:    "grp_1",
					ContractID: "ctr_C-123",
					Name:       propertyName2,
				}).Return(&propertyID2, nil)

			},
			expectedAttributes: map[string]string{
				"blocked_properties.#": "2",
				"blocked_properties.0": "prp_123456",
				"blocked_properties.1": "prp_456789",
			},
			expectError: nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(im *iam.Mock, _ *papi.Mock) {
				im.On("ListBlockedProperties", mock.Anything, iam.ListBlockedPropertiesRequest{
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
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			papiClient := &papi.Mock{}

			if test.init != nil {
				test.init(client, papiClient)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_blocked_properties.test", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_iam_blocked_properties.test", v))
			}
			useIAMandPAPIClient(client, papiClient, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataBlockedProperties/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
