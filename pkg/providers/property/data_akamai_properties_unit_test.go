package property

import (
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataProperties(t *testing.T) {
	t.Run("list properties", func(t *testing.T) {
		stateVal := compactJSON(loadFixtureBytes(fmt.Sprintf("testdata/TestDataProperties/%s", "properties.json")))
		client := &mockpapi{}
		props := papi.PropertiesItems{Items: []*papi.Property{
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note1",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp1",
				PropertyName:      "prpname1",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note2",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp2",
				PropertyName:      "prpname2",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
		}}

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "id", "grp_testctr_test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttrSet("data.akamai_properties.akaproperties", "properties"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "properties", stateVal),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list properties without group prefix", func(t *testing.T) {
		client := &mockpapi{}
		stateVal := compactJSON(loadFixtureBytes(fmt.Sprintf("testdata/TestDataProperties/%s", "properties.json")))
		props := papi.PropertiesItems{Items: []*papi.Property{
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note1",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp1",
				PropertyName:      "prpname1",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note2",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp2",
				PropertyName:      "prpname2",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
		}}

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties_no_group_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "id", "grp_testctr_test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "group_id", "test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttrSet("data.akamai_properties.akaproperties", "properties"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "properties", stateVal),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list properties without contract prefix", func(t *testing.T) {
		client := &mockpapi{}
		stateVal := compactJSON(loadFixtureBytes(fmt.Sprintf("testdata/TestDataProperties/%s", "properties.json")))
		props := papi.PropertiesItems{Items: []*papi.Property{
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note1",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp1",
				PropertyName:      "prpname1",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
			{
				AccountID:         "act1",
				AssetID:           "ast1",
				ContractID:        "ctr_test",
				GroupID:           "grp_test",
				LatestVersion:     1,
				Note:              "note2",
				ProductID:         "prd1",
				ProductionVersion: nil,
				PropertyID:        "prp2",
				PropertyName:      "prpname2",
				RuleFormat:        "latest",
				StagingVersion:    nil,
			},
		}}

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties_no_contract_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "id", "grp_testctr_test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "contract_id", "test"),
						resource.TestCheckResourceAttrSet("data.akamai_properties.akaproperties", "properties"),
						resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "properties", stateVal),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}
