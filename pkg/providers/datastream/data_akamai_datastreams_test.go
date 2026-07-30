package datastream

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	streamListWithOptionalIDs = []datastream.StreamDetails{
		{
			StreamID:      4,
			StreamName:    "Stream4",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       0,
			ContractID:    "",
			ProductID:     "P-1234",
			CreatedBy:     "user4",
			CreatedDate:   "02-08-2020 07:07:40 GMT",
			Properties: []datastream.Property{
				{
					PropertyID:   54375437,
					PropertyName: "property_name_5",
				},
			},
		},
	}

	streamListWithZeroGroupPopulatedContract = []datastream.StreamDetails{
		{
			StreamID:      5,
			StreamName:    "Stream5",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       0,
			ContractID:    "1-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user5",
			CreatedDate:   "03-08-2020 07:07:40 GMT",
			Properties: []datastream.Property{
				{
					PropertyID:   64376437,
					PropertyName: "property_name_6",
				},
			},
		},
	}

	streamListWithPopulatedGroupEmptyContract = []datastream.StreamDetails{
		{
			StreamID:      6,
			StreamName:    "Stream6",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       4321,
			ContractID:    "",
			ProductID:     "P-1234",
			CreatedBy:     "user6",
			CreatedDate:   "04-08-2020 07:07:40 GMT",
			Properties: []datastream.Property{
				{
					PropertyID:   74377437,
					PropertyName: "property_name_7",
				},
			},
		},
	}

	streamListForSpecificGroup = []datastream.StreamDetails{streamList[1]}

	appSecStreamList = []datastream.StreamDetails{
		{
			LogType:       datastream.LogTypeAppSec,
			StreamID:      10,
			StreamName:    "AppSecStream1",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       1234,
			ContractID:    "1-ABCDE",
			ProductID:     "KSD",
			CreatedBy:     "user1",
			CreatedDate:   "01-01-2024 00:00:00 GMT",
			AppSecConfigs: []datastream.AppSecConfig{
				{AppSecID: 12345, AppSecName: "WAF Security File"},
			},
		},
	}
)

func TestDataDatastreams(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		init  func(*datastream.Mock)
		steps []resource.TestStep
	}{
		"list streams": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(streamList, nil).Times(3)
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
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{
					GroupID: ptr.To(1234),
					LogType: datastream.LogTypeCDN, // default log type
				}).Return(streamListForSpecificGroup, nil).Times(3)
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
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{
					LogType: datastream.LogTypeCDN, // default log type
					GroupID: ptr.To(1234),
				}).Return(streamListForSpecificGroup, nil).Times(3)
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
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return([]datastream.StreamDetails{}, nil).Times(3)
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
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(streamListWithEmptyFields, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamListWithEmptyFields),
				},
			},
		},
		"list streams with optional group_id and contract_id": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(streamListWithOptionalIDs, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamListWithOptionalIDs),
				},
			},
		},
		"list streams with zero group_id and populated contract_id": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(streamListWithZeroGroupPopulatedContract, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamListWithZeroGroupPopulatedContract),
				},
			},
		},
		"list streams with populated group_id and empty contract_id": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(streamListWithPopulatedGroupEmptyContract, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					Check:  streamsChecks(streamListWithPopulatedGroupEmptyContract),
				},
			},
		},
		"could not fetch stream list": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{LogType: datastream.LogTypeCDN}).
					Return(nil, fmt.Errorf("failed to get stream list")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_without_groupid.tf"),
					ExpectError: regexp.MustCompile("failed to get stream list"),
				},
			},
		},
		"list appsec streams": {
			init: func(m *datastream.Mock) {
				// read
				m.On("ListStreams", testutils.MockContext, datastream.ListStreamsRequest{
					LogType: datastream.LogTypeAppSec,
				}).Return(appSecStreamList, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDatastreams/list_streams_appsec.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.#", "1"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.stream_id", "10"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.stream_name", "AppSecStream1"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.properties.#", "0"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.app_sec_configs.#", "1"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.app_sec_configs.0.app_sec_id", "12345"),
						resource.TestCheckResourceAttr("data.akamai_datastreams.test", "streams_details.0.app_sec_configs.0.app_sec_name", "WAF Security File"),
					),
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

func TestCreateStreamsAttrs_optionalGroupAndContractIDs(t *testing.T) {
	t.Parallel()

	attrs := createStreamsAttrs([]datastream.StreamDetails{
		{
			StreamID:      99,
			StreamName:    "optional-ids",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       0,
			ContractID:    "",
			ProductID:     "P-1234",
			CreatedBy:     "user",
			CreatedDate:   "01-01-2024 00:00:00 GMT",
		},
	})

	require.Len(t, attrs, 1)
	streamAttr, ok := attrs[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, 0, streamAttr["group_id"])
	assert.Equal(t, "", streamAttr["contract_id"])
	_, hasIntegrationType := streamAttr["integration_type"]
	assert.False(t, hasIntegrationType)
}

func TestCreateStreamsAttrs_omittedGroupWithPopulatedContract(t *testing.T) {
	t.Parallel()

	attrs := createStreamsAttrs([]datastream.StreamDetails{
		{
			StreamID:      100,
			StreamName:    "omitted-group",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       0,
			ContractID:    "1-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user",
			CreatedDate:   "01-01-2024 00:00:00 GMT",
		},
	})

	require.Len(t, attrs, 1)
	streamAttr := attrs[0].(map[string]interface{})
	assert.Equal(t, 0, streamAttr["group_id"])
	assert.Equal(t, "1-ABCDE", streamAttr["contract_id"])
}

func TestCreateStreamsAttrs_explicitZeroGroupWithPopulatedContract(t *testing.T) {
	t.Parallel()

	attrs := createStreamsAttrs([]datastream.StreamDetails{
		{
			StreamID:      102,
			StreamName:    "zero-group",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       0,
			ContractID:    "1-ABCDE",
			ProductID:     "P-1234",
			CreatedBy:     "user",
			CreatedDate:   "01-01-2024 00:00:00 GMT",
		},
	})

	require.Len(t, attrs, 1)
	streamAttr := attrs[0].(map[string]interface{})
	assert.Equal(t, 0, streamAttr["group_id"])
	assert.Equal(t, "1-ABCDE", streamAttr["contract_id"])
}

func TestCreateStreamsAttrs_populatedGroupWithEmptyContract(t *testing.T) {
	t.Parallel()

	attrs := createStreamsAttrs([]datastream.StreamDetails{
		{
			StreamID:      101,
			StreamName:    "empty-contract",
			StreamStatus:  datastream.StreamStatusActivated,
			StreamVersion: 1,
			LatestVersion: 1,
			GroupID:       4321,
			ContractID:    "",
			ProductID:     "P-1234",
			CreatedBy:     "user",
			CreatedDate:   "01-01-2024 00:00:00 GMT",
		},
	})

	require.Len(t, attrs, 1)
	streamAttr := attrs[0].(map[string]interface{})
	assert.Equal(t, 4321, streamAttr["group_id"])
	assert.Equal(t, "", streamAttr["contract_id"])
}

func TestCreateStreamsAttrs_setsIntegrationTypeOnlyWhenPresent(t *testing.T) {
	t.Parallel()

	withType := createStreamsAttrs([]datastream.StreamDetails{{
		StreamID: 1, StreamName: "a", StreamStatus: datastream.StreamStatusActivated,
		StreamVersion: 1, LatestVersion: 1, ProductID: "P", CreatedBy: "u", CreatedDate: "d",
		IntegrationType: "HYBRID",
	}})
	withoutType := createStreamsAttrs([]datastream.StreamDetails{{
		StreamID: 2, StreamName: "b", StreamStatus: datastream.StreamStatusActivated,
		StreamVersion: 1, LatestVersion: 1, ProductID: "P", CreatedBy: "u", CreatedDate: "d",
	}})

	_, hasWith := withType[0].(map[string]interface{})["integration_type"]
	_, hasWithout := withoutType[0].(map[string]interface{})["integration_type"]
	assert.True(t, hasWith)
	assert.False(t, hasWithout)
}

func TestCreateStreamsAttrs_alwaysEmitsContractAndGroupIDKeys(t *testing.T) {
	t.Parallel()

	attrs := createStreamsAttrs([]datastream.StreamDetails{{
		StreamID: 1, StreamName: "a", StreamStatus: datastream.StreamStatusActivated,
		StreamVersion: 1, LatestVersion: 1, GroupID: 0, ContractID: "",
		ProductID: "P", CreatedBy: "u", CreatedDate: "d",
	}})

	streamAttr := attrs[0].(map[string]interface{})
	_, hasGroupID := streamAttr["group_id"]
	_, hasContractID := streamAttr["contract_id"]
	assert.True(t, hasGroupID, "group_id key must always be present in datasource attrs")
	assert.True(t, hasContractID, "contract_id key must always be present in datasource attrs")
	assert.Equal(t, 0, streamAttr["group_id"])
	assert.Equal(t, "", streamAttr["contract_id"])
}
