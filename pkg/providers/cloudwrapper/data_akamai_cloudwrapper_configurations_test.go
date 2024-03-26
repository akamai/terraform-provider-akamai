package cloudwrapper

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataConfigurations(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudwrapper.Mock, []testDataForCWConfiguration)
		mockData   []testDataForCWConfiguration
		error      *regexp.Regexp
	}{
		"happy path- minimal data returned": {
			configPath: "testdata/TestDataConfigurations/default.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData []testDataForCWConfiguration) {
				expectGetConfigurations(m, testData, 5)
			},
			mockData: []testDataForCWConfiguration{
				minimalConfiguration,
			},
		},
		"happy path - all fields": {
			configPath: "testdata/TestDataConfigurations/default.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData []testDataForCWConfiguration) {
				expectGetConfigurations(m, testData, 5)
			},
			mockData: []testDataForCWConfiguration{
				configuration,
			},
		},
		"happy path - a few configurations": {
			configPath: "testdata/TestDataConfigurations/default.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData []testDataForCWConfiguration) {
				expectGetConfigurations(m, testData, 5)
			},
			mockData: []testDataForCWConfiguration{
				minimalConfiguration,
				configuration,
			},
		},
		"error getting configuration": {
			configPath: "testdata/TestDataConfigurations/default.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData []testDataForCWConfiguration) {
				m.On("ListConfigurations", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("get configuration failed")).Times(1)
			},
			mockData: []testDataForCWConfiguration{
				{
					ID: 1,
				},
			},
			error: regexp.MustCompile("get configuration failed"),
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

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
						Check:       checkCloudWrapperConfigurationsAttrs(test.mockData),
						ExpectError: test.error,
					},
				},
			})

			client.AssertExpectations(t)
		})
	}
}

func checkCloudWrapperConfigurationsAttrs(data []testDataForCWConfiguration) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_configurations.test", "configurations.#", strconv.Itoa(len(data))))
	for i, c := range data {
		checkFuncs = append(checkFuncs, checkCloudWrapperConfiguration(c, i))
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func expectGetConfigurations(client *cloudwrapper.Mock, data []testDataForCWConfiguration, timesToRun int) {
	var configurations []cloudwrapper.Configuration

	for _, c := range data {
		configurations = append(configurations, getConfiguration(c))
	}

	res := cloudwrapper.ListConfigurationsResponse{
		Configurations: configurations,
	}
	client.On("ListConfigurations", mock.Anything, mock.Anything).Return(&res, nil).Times(timesToRun)
}

func checkCloudWrapperConfiguration(data testDataForCWConfiguration, idx int) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	dsName := "data.akamai_cloudwrapper_configurations.test"
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, "id", "akamai_cloudwrapper_configurations"))
	checkFuncs = append(checkFuncs, checkConfiguration(data, dsName, "configurations."+strconv.Itoa(idx)+"."))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
