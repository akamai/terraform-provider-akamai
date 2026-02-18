package datastream

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var (
	streamList = []datastream.StreamDetails{
		{
			StreamID:        1,
			StreamName:      "Stream1",
			StreamStatus:    datastream.StreamStatusDeactivated,
			StreamVersion:   2,
			LatestVersion:   2,
			GroupID:         1234,
			ContractID:      "1-ABCDE",
			ProductID:       "P-1234",
			CreatedBy:       "user1",
			CreatedDate:     "14-07-2020 07:07:40 GMT",
			IntegrationType: "PM_DEPENDENT",
			Properties: []datastream.Property{
				{
					PropertyID:      13371337,
					PropertyName:    "property_name_1",
					IntegrationType: "PM_DEPENDENT",
				},
			},
		},
		{
			StreamID:        2,
			StreamName:      "Stream2",
			StreamStatus:    datastream.StreamStatusActivated,
			StreamVersion:   3,
			LatestVersion:   3,
			GroupID:         4321,
			ContractID:      "2-ABCDE",
			ProductID:       "P-1234",
			CreatedBy:       "user2",
			CreatedDate:     "24-07-2020 07:07:40 GMT",
			IntegrationType: "HYBRID",
			Properties: []datastream.Property{
				{
					PropertyID:      23372337,
					PropertyName:    "property_name_2",
					IntegrationType: "HYBRID",
				},
				{
					PropertyID:      33373337,
					PropertyName:    "property_name_3",
					IntegrationType: "DS_MANAGED",
				},
			},
		},
	}

	streamListWithEmptyFields = []datastream.StreamDetails{
		{
			StreamID:      3,
			StreamName:    "Stream3",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       1234,
			ContractID:    "1-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user3",
			CreatedDate:   "01-08-2020 07:07:40 GMT",
			// IntegrationType: "" (empty string when not set)
			Properties: []datastream.Property{
				{
					PropertyID:   44374437,
					PropertyName: "property_name_4",
					// IntegrationType: "" (empty string when not set)
				},
			},
		},
	}

	streamListForSpecificGroup = []datastream.StreamDetails{streamList[1]}
)

func TestDataDatastreams(t *testing.T) {
	tests := map[string]struct {
		init  func(*datastream.Mock)
		steps []resource.TestStep
	}{
		"list streams": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, mock.Anything).
					Return(streamList, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamList),
				},
			},
		},
		"list streams with specified group id": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{
					GroupID: ptr.To(1234),
				}).Return(streamListForSpecificGroup, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_with_groupid.tf"),
					Check:  streamsChecks(streamListForSpecificGroup),
				},
			},
		},
		"list streams with specified group id using grp prefix": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{
					GroupID: ptr.To(1234),
				}).Return(streamListForSpecificGroup, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_with_groupid_with_prefix.tf"),
					Check:  streamsChecks(streamListForSpecificGroup),
				},
			},
		},
		"list streams with specified group id using invalid prefix": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_with_groupid_with_invalid_prefix.tf"),
					ExpectError: regexp.MustCompile("Invalid reference"),
				},
			},
		},
		"list streams - empty list": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{}).
					Return([]datastream.StreamDetails{}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks([]datastream.StreamDetails{}),
				},
			},
		},
		"list streams with empty integration_type": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, mock.Anything).
					Return(streamListWithEmptyFields, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamListWithEmptyFields),
				},
			},
		},
		"could not fetch stream list": {
			init: func(m *datastream.Mock) {
				m.On("ListStreams", testutils.MockContext, mock.Anything).
					Return(nil, fmt.Errorf("failed to get stream list")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					ExpectError: regexp.MustCompile("failed to get stream list"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
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
		testCheck := []resource.TestCheckFunc{
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
		}
		// Only check integration_type if it's non-empty (API may not return it)
		if stream.IntegrationType != "" {
			testCheck = append(testCheck, resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("integration_type"), stream.IntegrationType))
		}
		checks = append(checks, resource.ComposeAggregateTestCheckFunc(testCheck...))
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
		// Only check integration_type if it's non-empty (API may not return it)
		if property.IntegrationType != "" {
			testCheck = append(testCheck, resource.TestCheckResourceAttr("data.akamai_datastreams.test", attrName("integration_type"), property.IntegrationType))
		}
		checks = append(checks, testCheck...)
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}
