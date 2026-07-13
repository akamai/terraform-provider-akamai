package datastream

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceStreamLogTypeValidation(t *testing.T) {
	tests := map[string]struct {
		tfFile    string
		withError *regexp.Regexp
	}{
		"appsec missing app_sec_configs": {
			tfFile:    "testdata/TestResourceStream/appsec/appsec_missing_app_sec_configs.tf",
			withError: regexp.MustCompile("app_sec_configs.*required"),
		},
		"appsec with properties": {
			tfFile:    "testdata/TestResourceStream/appsec/appsec_with_properties.tf",
			withError: regexp.MustCompile("cannot set `properties` when log_type is \"APPSEC\""),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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

func TestResourceStreamAppSecCreateRead(t *testing.T) {
	client := &datastream.Mock{}

	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: 30,
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
		DatasetFields: []datastream.DatasetFieldID{},
		ContractID:    "test_contract",
		GroupID:       42,
		AppSecConfigs: []datastream.AppSecConfigID{
			{AppSecID: 16536},
		},
		StreamName:         "test-app-sec-stream-create",
		NotificationEmails: []string{"nobody@akamai.com"},
	}

	createReq := datastream.CreateStreamRequest{
		StreamConfiguration: streamConfiguration,
		Activate:            false,
		LogType:             datastream.LogTypeAppSec,
	}

	getReq := datastream.GetStreamRequest{
		StreamID: streamID,
		LogType:  datastream.LogTypeAppSec,
	}

	client.On("CreateStream", testutils.MockContext, createReq).
		Return(&datastream.DetailedStreamVersion{
			StreamID:      streamID,
			StreamVersion: 1,
			GroupID:       streamConfiguration.GroupID,
		}, nil).Once()

	client.On("GetStream", testutils.MockContext, getReq).
		Return(&datastream.DetailedStreamVersion{
			LogType:               getReq.LogType,
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
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
			AppSecConfigs: []datastream.AppSecConfig{
				{
					AppSecID:   16536,
					AppSecName: "WAF Security File",
				},
			},
			StreamID:           streamID,
			StreamName:         streamConfiguration.StreamName,
			StreamVersion:      1,
			LatestVersion:      1,
			NotificationEmails: streamConfiguration.NotificationEmails,
			ModifiedDate:       "01-01-2020 12:00:00 GMT",
		}, nil).Times(3)

	// for after the update step, where we expect the same GetStream call with the same response but with the updated email.
	client.On("GetStream", testutils.MockContext, getReq).
		Return(&datastream.DetailedStreamVersion{
			LogType:               getReq.LogType,
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
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
			AppSecConfigs: []datastream.AppSecConfig{
				{
					AppSecID:   16536,
					AppSecName: "WAF Security File",
				},
			},
			StreamID:           streamID,
			StreamName:         streamConfiguration.StreamName,
			StreamVersion:      1,
			LatestVersion:      1,
			ModifiedDate:       "01-01-2020 12:00:00 GMT",
			NotificationEmails: []string{"updated@akamai.com"},
		}, nil)

	client.On("DeleteStream", testutils.MockContext, datastream.DeleteStreamRequest{
		StreamID: streamID,
		LogType:  datastream.LogTypeAppSec,
	}).Return(' ', nil).Once()

	// everything for the update should be the same except the email address.
	updateStreamRequest := datastream.UpdateStreamRequest{
		StreamID:            streamID,
		LogType:             datastream.LogTypeAppSec,
		StreamConfiguration: streamConfiguration,
	}
	updateStreamRequest.StreamConfiguration.NotificationEmails = []string{"updated@akamai.com"}
	updateStreamRequest.StreamConfiguration.GroupID = 0 // group ID does not get sent for an update.

	client.On("UpdateStream", testutils.MockContext, updateStreamRequest).Return(&datastream.DetailedStreamVersion{
		LogType:               getReq.LogType,
		StreamStatus:          datastream.StreamStatusInactive,
		DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
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
		AppSecConfigs: []datastream.AppSecConfig{
			{
				AppSecID:   16536,
				AppSecName: "WAF Security File",
			},
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
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/appsec/create_app_sec_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
						resource.TestCheckResourceAttr("akamai_datastream.s", "log_type", string(datastream.LogTypeAppSec)),
						resource.TestCheckResourceAttr("akamai_datastream.s", "active", "false"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "properties.#", "0"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "app_sec_configs.#", "1"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "app_sec_configs.0", "16536"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/appsec/update_app_sec_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
						resource.TestCheckResourceAttr("akamai_datastream.s", "log_type", string(datastream.LogTypeAppSec)),
						resource.TestCheckResourceAttr("akamai_datastream.s", "active", "false"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "properties.#", "0"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "app_sec_configs.#", "1"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "app_sec_configs.0", "16536"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "1"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.0", "updated@akamai.com"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
