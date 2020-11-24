package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func TestDSPropertyRulesRead(t *testing.T) {
	t.Run("get datasource property rules", func(t *testing.T) {
		client := &mockpapi{}
		mockImpl := func(m *mockpapi) {
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(&papi.GetPropertyVersionsResponse{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				Version: papi.PropertyVersionGetItem{
					PropertyVersion: 1,
				},
			}, nil).Once()
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(&papi.GetPropertyVersionsResponse{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				Version: papi.PropertyVersionGetItem{
					PropertyVersion: 1,
				},
			}, nil).Once()
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
			}).Return(&papi.GetRuleTreeResponse{}, nil)
		}
		mockImpl(client)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						ExpectNonEmptyPlan: true,
						Config:             loadFixtureString("testdata/TestDSPropertyRules/ds_property_rules.tf"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestDSPropertyRulesRead_Fail(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config:      loadFixtureString("testdata/TestDSPropertyRules/always_fails.tf"),
			ExpectError: regexp.MustCompile(`Error: required value cannot be blank`),
		}},
	})
}
