package property

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
)

func TestDataProperties(t *testing.T) {
	t.Run("list properties", func(t *testing.T) {
		client := &papi.Mock{}
		props := papi.PropertiesItems{Items: buildPapiProperties()}
		properties := decodePropertyItems(props.Items)

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties.tf"),
					Check:  buildAggregatedTest(properties, "grp_testctr_test", "grp_test", "ctr_test"),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list properties without group prefix", func(t *testing.T) {
		client := &papi.Mock{}
		props := papi.PropertiesItems{Items: buildPapiProperties()}
		properties := decodePropertyItems(props.Items)

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties_no_group_prefix.tf"),
					Check:  buildAggregatedTest(properties, "grp_testctr_test", "test", "ctr_test"),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list properties without contract prefix", func(t *testing.T) {
		client := &papi.Mock{}
		props := papi.PropertiesItems{Items: buildPapiProperties()}
		properties := decodePropertyItems(props.Items)

		client.On("GetProperties",
			mock.Anything,
			papi.GetPropertiesRequest{GroupID: "grp_test", ContractID: "ctr_test"},
		).Return(&papi.GetPropertiesResponse{Properties: props}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDataProperties/properties_no_contract_prefix.tf"),
					Check:  buildAggregatedTest(properties, "grp_testctr_test", "grp_test", "test"),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}

func buildPapiProperties() []*papi.Property {
	properties := make([]*papi.Property, 10)
	for i := 0; i < 10; i++ {
		properties[i] = &papi.Property{
			AccountID:         fmt.Sprintf("act%v", i),
			AssetID:           fmt.Sprintf("ast%v", i),
			ContractID:        "ctr_test",
			GroupID:           "grp_test",
			LatestVersion:     1,
			Note:              fmt.Sprintf("note%v", i),
			ProductID:         "prd1",
			ProductionVersion: nil,
			PropertyID:        fmt.Sprintf("prp%v", i),
			PropertyName:      fmt.Sprintf("prpname%v", i),
			RuleFormat:        "latest",
			StagingVersion:    nil,
		}
	}
	return properties
}

func buildAggregatedTest(properties []map[string]interface{}, id, groupID, contractID string) resource.TestCheckFunc {
	testVar := make([]resource.TestCheckFunc, 0)
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "id", id))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "group_id", groupID))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "contract_id", contractID))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", "properties.#", fmt.Sprintf("%v", len(properties))))
	nrEntries := fmt.Sprintf("%v", len(properties[0]))
	for ind, property := range properties {
		keyNrEntries := fmt.Sprintf("properties.%v.%%", ind)
		testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", keyNrEntries, nrEntries))
		for mapKey, mapVal := range property {
			value := fmt.Sprintf("%v", mapVal)
			key := fmt.Sprintf("properties.%v.%v", ind, mapKey)
			testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_properties.akaproperties", key, value))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(testVar...)
}

func decodePropertyItems(items []*papi.Property) []map[string]interface{} {
	properties := make([]map[string]interface{}, 0)
	for _, item := range items {
		prop := map[string]interface{}{
			"contract_id":        item.ContractID,
			"group_id":           item.GroupID,
			"latest_version":     item.LatestVersion,
			"note":               item.Note,
			"product_id":         item.ProductID,
			"production_version": decodeVersion(item.ProductionVersion),
			"property_id":        item.PropertyID,
			"property_name":      item.PropertyName,
			"rule_format":        item.RuleFormat,
			"staging_version":    decodeVersion(item.StagingVersion),
		}
		properties = append(properties, prop)
	}
	return properties
}
