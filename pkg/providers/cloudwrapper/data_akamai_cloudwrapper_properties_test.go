package cloudwrapper

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataProperty(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudwrapper.Mock, testDataForCWProperties)
		mockData   testDataForCWProperties
		error      *regexp.Regexp
	}{
		"happy path - one property, unused-true": {
			configPath: "testdata/TestDataProperties/default_unused_true.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWProperties) {
				expectListProperties(m, testData, 5)
			},
			mockData: testDataForCWProperties{
				unused: true,
				properties: []cloudwrapper.Property{
					{
						GroupID:      1,
						ContractID:   "ctr_1",
						PropertyID:   11,
						PropertyName: "Name1",
						Type:         "Type1",
					},
				},
			},
		},
		"happy path - two properties, unused-false, contract_ids supplied": {
			configPath: "testdata/TestDataProperties/default_unused_false.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWProperties) {
				expectListProperties(m, testData, 5)
			},
			mockData: testDataForCWProperties{
				contractIDs: []string{"ctr_1", "ctr_2"},
				properties: []cloudwrapper.Property{
					{
						GroupID:      1,
						ContractID:   "ctr_1",
						PropertyID:   11,
						PropertyName: "Name1",
						Type:         "Type1",
					},
					{
						GroupID:      2,
						ContractID:   "ctr_2",
						PropertyID:   22,
						PropertyName: "Name2",
						Type:         "Type2",
					},
				},
			},
		},
		"happy path - no optional attributes": {
			configPath: "testdata/TestDataProperties/no_attributes.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWProperties) {
				expectListProperties(m, testData, 5)
			},
			mockData: testDataForCWProperties{
				properties: []cloudwrapper.Property{
					{
						GroupID:      1,
						ContractID:   "ctr_1",
						PropertyID:   11,
						PropertyName: "Name1",
						Type:         "Type1",
					},
				},
			},
		},
		"happy path - empty properties list": {
			configPath: "testdata/TestDataProperties/default_unused_false.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWProperties) {
				expectListProperties(m, testData, 5)
			},
			mockData: testDataForCWProperties{
				contractIDs: []string{"ctr_1", "ctr_2"},
				properties:  []cloudwrapper.Property{},
			},
		},
		"error listing properties": {
			configPath: "testdata/TestDataProperties/default_unused_false.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWProperties) {
				expectListPropertiesWithError(m, testData, 1)
			},
			mockData: testDataForCWProperties{
				contractIDs: []string{"ctr_1", "ctr_2"},
			},
			error: regexp.MustCompile("list properties failed"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudwrapper.Mock{}
			if test.init != nil {
				test.init(t, client, test.mockData)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, test.configPath),
						Check:       checkCloudWrapperPropertiesAttrs(test.mockData),
						ExpectError: test.error,
					},
				},
			})

			client.AssertExpectations(t)
		})
	}
}

type testDataForCWProperties struct {
	unused      bool
	contractIDs []string
	properties  []cloudwrapper.Property
}

func checkCloudWrapperPropertiesAttrs(data testDataForCWProperties) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", "contract_ids.#", strconv.Itoa(len(data.contractIDs))))
	for i, ctr := range data.contractIDs {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("contract_ids.%d", i), ctr))
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", "properties.#", strconv.Itoa(len(data.properties))))
	for i, prp := range data.properties {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("properties.%d.property_id", i), strconv.FormatInt(prp.PropertyID, 10)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("properties.%d.type", i), string(prp.Type)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("properties.%d.property_name", i), prp.PropertyName))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("properties.%d.contract_id", i), prp.ContractID))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_properties.test", fmt.Sprintf("properties.%d.group_id", i), strconv.FormatInt(prp.GroupID, 10)))
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttrSet("data.akamai_cloudwrapper_properties.test", "id"))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func expectListProperties(client *cloudwrapper.Mock, data testDataForCWProperties, timesToRun int) {
	listPropertiesReq := cloudwrapper.ListPropertiesRequest{
		Unused:      data.unused,
		ContractIDs: data.contractIDs,
	}
	listPropertiesRes := cloudwrapper.ListPropertiesResponse{
		Properties: data.properties,
	}
	client.On("ListProperties", mock.Anything, listPropertiesReq).Return(&listPropertiesRes, nil).Times(timesToRun)
}

func expectListPropertiesWithError(client *cloudwrapper.Mock, data testDataForCWProperties, timesToRun int) {
	listPropertiesReq := cloudwrapper.ListPropertiesRequest{
		Unused:      data.unused,
		ContractIDs: data.contractIDs,
	}
	client.On("ListProperties", mock.Anything, listPropertiesReq).Return(nil, fmt.Errorf("list properties failed")).Times(timesToRun)
}
