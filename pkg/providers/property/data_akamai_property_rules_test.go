package property

import (
	"fmt"
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
			}, nil)
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
			}).Return(&papi.GetRuleTreeResponse{
				Rules: papi.Rules{
					Name: "some rule tree",
				},
				Response: papi.Response{
					Errors: []*papi.Error{
						{
							Title: "some error",
						},
					},
				},
			}, nil)
		}
		mockImpl(client)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSPropertyRules/ds_property_rules.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "id", "prp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "property_id", "prp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "group_id", "grp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "contract_id", "ctr_2"),
							resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "rules"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "errors", `[{"type":"","title":"some error","detail":""}]`),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("group_id is required with contract_id", func(t *testing.T) {
		client := &mockpapi{}
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/missing_group_id.tf"),
						ExpectError: regexp.MustCompile("\"contract_id\": all of `contract_id,group_id` must be specified"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("contract_id is required with group_id", func(t *testing.T) {
		client := &mockpapi{}
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/missing_contract_id.tf"),
						ExpectError: regexp.MustCompile("\"group_id\": all of `contract_id,group_id` must be specified"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("contract_id cannot be empty", func(t *testing.T) {
		client := &mockpapi{}
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/empty_contract_id.tf"),
						ExpectError: regexp.MustCompile(`provided value cannot be blank((.|\n)*)contract_id = ""`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("group_id cannot be empty", func(t *testing.T) {
		client := &mockpapi{}
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/empty_group_id.tf"),
						ExpectError: regexp.MustCompile(`provided value cannot be blank((.|\n)*)group_id = ""`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error fetching latest version", func(t *testing.T) {
		client := &mockpapi{}
		mockImpl := func(m *mockpapi) {
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(nil, fmt.Errorf("fetching latest version")).Once()
		}
		mockImpl(client)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/ds_property_rules.tf"),
						ExpectError: regexp.MustCompile("fetching latest version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error fetching rules", func(t *testing.T) {
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
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
			}).Return(nil, fmt.Errorf("fetching rule tree")).Once()
		}
		mockImpl(client)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSPropertyRules/ds_property_rules.tf"),
						ExpectError: regexp.MustCompile("fetching rule tree"),
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
			ExpectError: regexp.MustCompile(`Error: provided value cannot be blank`),
		}},
	})
}
