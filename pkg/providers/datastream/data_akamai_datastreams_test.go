package datastream

import (
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/stretchr/testify/mock"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	streamList = []datastream.StreamDetails{
		{
			StreamID:      1,
			StreamName:    "Stream1",
			StreamStatus:  datastream.StreamStatusDeactivated,
			StreamVersion: 2,
			LatestVersion: 2,
			GroupID:       1234,
			ContractID:    "1-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user1",
			CreatedDate:   "14-07-2020 07:07:40 GMT",
			Properties: []datastream.Property{
				{
					PropertyID:   13371337,
					PropertyName: "property_name_1",
				},
			},
		},
		{
			StreamID:      2,
			StreamName:    "Stream2",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 3,
			LatestVersion: 3,
			GroupID:       4321,
			ContractID:    "2-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user2",
			CreatedDate:   "24-07-2020 07:07:40 GMT",
			Properties: []datastream.Property{
				{
					PropertyID:   23372337,
					PropertyName: "property_name_2",
				},
				{
					PropertyID:   33373337,
					PropertyName: "property_name_3",
				},
			},
		},
	}

	streamListForSpecificGroup = []datastream.StreamDetails{streamList[1]}
)

func TestDataDatastreams(t *testing.T) {
	tests := map[string]struct {
		init  func(*testing.T, *datastream.Mock)
		steps []resource.TestStep
	}{
		"list streams": {
			init: func(t *testing.T, m *datastream.Mock) {
				m.On("ListStreams", mock.Anything, mock.Anything).
					Return(streamList, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamList),
				},
			},
		},
		"list streams with specified group id": {
			init: func(t *testing.T, m *datastream.Mock) {
				m.On("ListStreams", mock.Anything, datastream.ListStreamsRequest{
					GroupID: tools.IntPtr(1234),
				}).Return(streamListForSpecificGroup, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataDatastreams/list_streams_with_groupid.tf"),
					Check:  streamsChecks(streamListForSpecificGroup),
				},
			},
		},
		"list streams with specified group id using grp prefix": {
			init: func(t *testing.T, m *datastream.Mock) {
				m.On("ListStreams", mock.Anything, datastream.ListStreamsRequest{
					GroupID: tools.IntPtr(1234),
				}).Return(streamListForSpecificGroup, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataDatastreams/list_streams_with_groupid_with_prefix.tf"),
					Check:  streamsChecks(streamListForSpecificGroup),
				},
			},
		},
		"list streams with specified group id using invalid prefix": {
			init: func(t *testing.T, m *datastream.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("testdata/TestDataDatastreams/list_streams_with_groupid_with_invalid_prefix.tf"),
					ExpectError: regexp.MustCompile("Invalid reference"),
				},
			},
		},
		"list streams - empty list": {
			init: func(t *testing.T, m *datastream.Mock) {
				m.On("ListStreams", mock.Anything, datastream.ListStreamsRequest{}).
					Return([]datastream.StreamDetails{}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks([]datastream.StreamDetails{}),
				},
			},
		},
		"could not fetch stream list": {
			init: func(t *testing.T, m *datastream.Mock) {
				m.On("ListStreams", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("failed to get stream list")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					ExpectError: regexp.MustCompile("failed to get stream list"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}
			test.init(t, client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps:             test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func streamsChecks(details []datastream.StreamDetails) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.#", strconv.Itoa(len(details))),
	}
	for idx, stream := range details {
		attrName := func(attr string) string { return fmt.Sprintf("streams_details.%d.%s", idx, attr) }
		testCheck := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("stream_id"), strconv.FormatInt(stream.StreamID, 10)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("stream_name"), stream.StreamName),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("stream_version"), strconv.FormatInt(stream.StreamVersion, 10)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("group_id"), strconv.Itoa(stream.GroupID)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("contract_id"), stream.ContractID),
			propertiesCheck(attrName("properties"), stream.Properties),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("latest_version"), strconv.FormatInt(stream.LatestVersion, 10)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("product_id"), stream.ProductID),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("stream_status"), string(stream.StreamStatus)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("created_by"), stream.CreatedBy),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("created_date"), stream.CreatedDate),
		)
		checks = append(checks, testCheck)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func propertiesCheck(key string, properties []datastream.Property) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.akamai_datastreams.test", fmt.Sprintf("%s.#", key), strconv.Itoa(len(properties))),
	}
	for idx, property := range properties {
		attrName := func(attr string) string { return fmt.Sprintf("%s.%d.%s", key, idx, attr) }
		testCheck := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("property_id"), strconv.Itoa(property.PropertyID)),
			resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("property_name"), property.PropertyName),
		}
		checks = append(checks, testCheck...)
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}
