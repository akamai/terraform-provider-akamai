package datastream

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/datastream"
	test "github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceStreamAnswerXLogTypeValidation(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		tfFile    string
		withError *regexp.Regexp
	}{
		"answerx missing dataset_fields": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_missing_dataset_fields.tf",
			withError: regexp.MustCompile("`dataset_fields` are required for log_type \"ANSWERX\""),
		},
		"answerx empty dataset_fields": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_empty_dataset_fields.tf",
			withError: regexp.MustCompile("`dataset_fields` are required for log_type \"ANSWERX\""),
		},
		"answerx with properties": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_with_properties.tf",
			withError: regexp.MustCompile("cannot set `properties` when log_type is \"ANSWERX\""),
		},
		"answerx with app_sec_configs": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_with_app_sec_configs.tf",
			withError: regexp.MustCompile("cannot set `app_sec_configs` when log_type is \"ANSWERX\""),
		},
		"answerx missing service_ids": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_missing_service_ids.tf",
			withError: regexp.MustCompile("`service_ids` are required for log_type \"ANSWERX\""),
		},
		"answerx empty service_ids": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_empty_service_ids.tf",
			withError: regexp.MustCompile("`service_ids` are required for log_type \"ANSWERX\""),
		},
		"answerx negative service_ids": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_negative_service_ids.tf",
			withError: regexp.MustCompile(`at least \(0\), got -1`),
		},
		"answerx missing contract_id": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_missing_contract_id.tf",
			withError: regexp.MustCompile("`contract_id` is required for log_type \"ANSWERX\""),
		},
		"answerx whitespace contract_id": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_whitespace_contract_id.tf",
			withError: regexp.MustCompile("`contract_id` is required for log_type \"ANSWERX\""),
		},
		"answerx missing group_id": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_missing_group_id.tf",
			withError: regexp.MustCompile("`group_id` is required for log_type \"ANSWERX\""),
		},
		"answerx invalid group_id": {
			tfFile:    "testdata/TestResourceStream/answerx/answerx_invalid_group_id.tf",
			withError: regexp.MustCompile("invalid `group_id`"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &datastream.Mock{}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.tfFile),
							ExpectError: test.withError,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func TestResourceStreamLogTypeValidationUnknownCollectionsPlanOnly(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"answerx unknown required fields are deferred":     "testdata/TestResourceStream/unknown_plan/answerx_unknown_dataset_fields_and_service_ids.tf",
		"cdn unknown forbidden service_ids is deferred":    "testdata/TestResourceStream/unknown_plan/cdn_unknown_forbidden_service_ids.tf",
		"appsec unknown forbidden service_ids is deferred": "testdata/TestResourceStream/unknown_plan/appsec_unknown_forbidden_service_ids.tf",
	}

	for name, tfFile := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &datastream.Mock{}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:             testutils.LoadFixtureString(t, tfFile),
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func TestResourceStreamAnswerXCreateRead(t *testing.T) {
	t.Parallel()
	client := &datastream.Mock{}

	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
		},
		Destination: datastream.AbstractConnector(
			&datastream.TrafficPeakConnector{
				ContentType:        "application/json",
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				DisplayName:        "TrafficPeakTest",
				Endpoint:           "https://example.com/ingest/event?table=unit_test&token=1234",
				UserName:           "username",
				Password:           "password",
			},
		),
		DatasetFields: []datastream.DatasetFieldID{
			{DatasetFieldID: 2000},
		},
		ContractID: "test_contract",
		GroupID:    42,
		AnswerXServiceIDs: []datastream.AnswerXServiceID{
			{SSID: 101},
		},
		StreamName:         "test-answerx-stream",
		NotificationEmails: []string{"abc@akamai.com"},
	}

	createReq := datastream.CreateStreamRequest{
		StreamConfiguration: streamConfiguration,
		Activate:            false,
		LogType:             datastream.LogTypeAnswerX,
	}

	getReq := datastream.GetStreamRequest{
		StreamID: streamID,
		LogType:  datastream.LogTypeAnswerX,
	}

	// API response populates name/product fields in AnswerXServiceIDs; provider reads only the SSID.
	commonGetResp := &datastream.DetailedStreamVersion{
		LogType:      getReq.LogType,
		StreamStatus: datastream.StreamStatusInactive,
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
		},
		Destination: datastream.Destination{
			DestinationType:    datastream.DestinationTypeTrafficPeak,
			AuthenticationType: datastream.AuthenticationTypeBasic,
			CompressLogs:       true,
			DisplayName:        "TrafficPeakTest",
			Endpoint:           "https://example.com/ingest/event?table=unit_test&token=1234",
			ContentType:        "application/json",
		},
		ContractID: streamConfiguration.ContractID,
		GroupID:    streamConfiguration.GroupID,
		DatasetFields: []datastream.DataSetField{
			{DatasetFieldID: 2000},
		},
		AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
			{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
		},
		StreamID:           streamID,
		StreamName:         streamConfiguration.StreamName,
		StreamVersion:      1,
		LatestVersion:      1,
		NotificationEmails: streamConfiguration.NotificationEmails,
		ModifiedDate:       "01-01-2020 12:00:00 GMT",
	}

	client.On("CreateStream", testutils.MockContext, createReq).
		Return(&datastream.DetailedStreamVersion{
			StreamID:      streamID,
			StreamVersion: 1,
			GroupID:       streamConfiguration.GroupID,
		}, nil).Once()

	// Set up GetStream responses: the first responses return the original data, and later ones return updated data.
	getRespUpdated := *commonGetResp
	getRespUpdated.NotificationEmails = []string{"updated@akamai.com"}
	getRespUpdated.AnswerXServiceIDs = []datastream.AnswerXServiceDetail{
		{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
		{SSID: 102, Name: "ServiceB", Product: "AnswerX"},
	}

	// First batch of GetStream calls return original response
	// This covers: post-create Read, post-apply plan refresh, pre-update status check
	client.On("GetStream", testutils.MockContext, getReq).
		Return(commonGetResp, nil).Times(3)

	// Use the updated response for the pre-update status check, update status polling, post-update Read, and delete status check.
	client.On("GetStream", testutils.MockContext, getReq).
		Return(&getRespUpdated, nil).Times(4)

	client.On("DeleteStream", testutils.MockContext, datastream.DeleteStreamRequest{
		StreamID: streamID,
		LogType:  datastream.LogTypeAnswerX,
	}).Return(nil).Once()

	// At update time, state has service_ids with the server-returned name/product values.
	updateStreamRequest := datastream.UpdateStreamRequest{
		StreamID: streamID,
		StreamConfiguration: datastream.StreamConfiguration{
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination:           streamConfiguration.Destination,
			ContractID:            streamConfiguration.ContractID,
			DatasetFields: []datastream.DatasetFieldID{
				{DatasetFieldID: 2000},
			},
			AnswerXServiceIDs: []datastream.AnswerXServiceID{
				{SSID: 101},
				{SSID: 102},
			},
			StreamName:         streamConfiguration.StreamName,
			NotificationEmails: []string{"updated@akamai.com"},
		},
		Activate: false,
		LogType:  datastream.LogTypeAnswerX,
	}

	client.On("UpdateStream", testutils.MockContext, mock.MatchedBy(func(req datastream.UpdateStreamRequest) bool {
		if req.StreamID != streamID || req.Activate != updateStreamRequest.Activate || req.LogType != updateStreamRequest.LogType {
			return false
		}
		cfg := req.StreamConfiguration
		if cfg.ContractID != updateStreamRequest.StreamConfiguration.ContractID || cfg.GroupID != updateStreamRequest.StreamConfiguration.GroupID || cfg.StreamName != updateStreamRequest.StreamConfiguration.StreamName {
			return false
		}
		if !sameAnswerXServiceIDSet(cfg.AnswerXServiceIDs, updateStreamRequest.StreamConfiguration.AnswerXServiceIDs) {
			return false
		}
		if len(cfg.NotificationEmails) != len(updateStreamRequest.StreamConfiguration.NotificationEmails) {
			return false
		}
		for i := range cfg.NotificationEmails {
			if cfg.NotificationEmails[i] != updateStreamRequest.StreamConfiguration.NotificationEmails[i] {
				return false
			}
		}
		return true
	})).Return(&datastream.DetailedStreamVersion{
		LogType:      getReq.LogType,
		StreamStatus: datastream.StreamStatusInactive,
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
		},
		Destination: datastream.Destination{
			DestinationType:    datastream.DestinationTypeTrafficPeak,
			AuthenticationType: datastream.AuthenticationTypeBasic,
			CompressLogs:       true,
			DisplayName:        "TrafficPeakTest",
			Endpoint:           "https://example.com/ingest/event?table=unit_test&token=1234",
			ContentType:        "application/json",
		},
		ContractID: streamConfiguration.ContractID,
		GroupID:    streamConfiguration.GroupID,
		DatasetFields: []datastream.DataSetField{
			{DatasetFieldID: 2000},
		},
		AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
			{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
			{SSID: 102, Name: "ServiceB", Product: "AnswerX"},
		},
		StreamID:           streamID,
		StreamName:         streamConfiguration.StreamName,
		StreamVersion:      1,
		LatestVersion:      1,
		ModifiedDate:       "01-01-2020 12:00:00 GMT",
		NotificationEmails: []string{"updated@akamai.com"},
	}, nil).Once()

	useClient(client, func() {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					//Create
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/answerx/create_answerx_stream.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						test.NewStateChecker("akamai_datastream.s").
							CheckEqual("id", strconv.FormatInt(streamID, 10)).
							CheckEqual("log_type", string(datastream.LogTypeAnswerX)).
							CheckEqual("active", "false").
							CheckEqual("properties.#", "0").
							CheckEqual("app_sec_configs.#", "0").
							CheckEqual("service_ids.#", "1").
							CheckEqual("dataset_fields.#", "1").
							CheckEqual("dataset_fields.0", "2000").
							Build(),
						resource.TestCheckTypeSetElemAttr("akamai_datastream.s", "service_ids.*", "101"),
					),
				},
				{
					//Update
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/answerx/update_answerx_stream.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						test.NewStateChecker("akamai_datastream.s").
							CheckEqual("id", strconv.FormatInt(streamID, 10)).
							CheckEqual("log_type", string(datastream.LogTypeAnswerX)).
							CheckEqual("active", "false").
							CheckEqual("properties.#", "0").
							CheckEqual("app_sec_configs.#", "0").
							CheckEqual("service_ids.#", "2").
							CheckEqual("notification_emails.#", "1").
							CheckEqual("notification_emails.0", "updated@akamai.com").
							Build(),
						resource.TestCheckTypeSetElemAttr("akamai_datastream.s", "service_ids.*", "101"),
						resource.TestCheckTypeSetElemAttr("akamai_datastream.s", "service_ids.*", "102"),
					),
				},
			},
		})
	})

	client.AssertExpectations(t)
}
