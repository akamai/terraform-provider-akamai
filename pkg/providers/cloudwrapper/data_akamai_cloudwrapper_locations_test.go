package cloudwrapper

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataLocations(t *testing.T) {
	expectListLocations := func(client *cloudwrapper.Mock, data testDataForCWLocations, timesToRun int) {
		listLocationsRes := cloudwrapper.ListLocationResponse{
			Locations: data.locations,
		}
		client.On("ListLocations", testutils.MockContext).Return(&listLocationsRes, nil).Times(timesToRun)
	}

	expectListLocationsWithError := func(client *cloudwrapper.Mock, timesToRun int) {
		client.On("ListLocations", testutils.MockContext).Return(nil, fmt.Errorf("list locations failed")).Times(timesToRun)
	}

	location := testDataForCWLocations{
		locations: []cloudwrapper.Location{
			{
				LocationID:   1,
				LocationName: "US East",
				TrafficTypes: []cloudwrapper.TrafficTypeItem{
					{
						TrafficTypeID: 1,
						TrafficType:   "LIVE",
						MapName:       "cw-s-use-live",
					},
					{
						TrafficTypeID: 2,
						TrafficType:   "LIVE_VOD",
						MapName:       "cw-s-use",
					},
				},
				MultiCDNLocationID: "018",
			},
			{
				LocationID:   2,
				LocationName: "US West",
				TrafficTypes: []cloudwrapper.TrafficTypeItem{
					{
						TrafficTypeID: 3,
						TrafficType:   "LIVE_VOD",
						MapName:       "cw-s-usw",
					},
					{
						TrafficTypeID: 4,
						TrafficType:   "LIVE",
						MapName:       "cw-s-usw-live",
					},
				},
				MultiCDNLocationID: "020",
			},
		},
	}
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudwrapper.Mock, testDataForCWLocations)
		mockData   testDataForCWLocations
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataLocations/location.tf",
			init: func(_ *testing.T, m *cloudwrapper.Mock, testData testDataForCWLocations) {
				expectListLocations(m, testData, 3)
			},
			mockData: location,
		},
		"error listing locations": {
			configPath: "testdata/TestDataLocations/location.tf",
			init: func(_ *testing.T, m *cloudwrapper.Mock, _ testDataForCWLocations) {
				expectListLocationsWithError(m, 1)
			},
			error: regexp.MustCompile("list locations failed"),
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
						Check:       checkCloudWrapperLocationsAttrs(),
						ExpectError: test.error,
					},
				},
			})

			client.AssertExpectations(t)
		})
	}
}

type testDataForCWLocations struct {
	locations []cloudwrapper.Location
}

func checkCloudWrapperLocationsAttrs() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_locations.test", "locations.#", "2"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_locations.test", "locations.1.location_id", "2"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_locations.test", "locations.1.traffic_types.#", "2"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_locations.test", "locations.1.traffic_types.0.traffic_type_id", "3"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_locations.test", "locations.1.traffic_types.0.location_id", "cw-s-usw"))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
