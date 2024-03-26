package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

const (
	streamID = int64(12321)
)

func TestResourceStream(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		client := &datastream.Mock{}

		PollForActivationStatusChangeInterval = 1 * time.Millisecond

		streamConfiguration := datastream.StreamConfiguration{
			CollectMidgress: false,
			DeliveryConfiguration: datastream.DeliveryConfiguration{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Destination: datastream.AbstractConnector(
				&datastream.S3Connector{
					AccessKey:       "s3_test_access_key",
					Bucket:          "s3_test_bucket",
					DisplayName:     "s3_test_connector_name",
					Path:            "s3_test_path",
					Region:          "s3_test_region",
					SecretAccessKey: "s3_test_secret_key",
				},
			),
			ContractID: "test_contract",
			DatasetFields: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				},
				{
					DatasetFieldID: 1002,
				},
				{
					DatasetFieldID: 2000,
				},
				{
					DatasetFieldID: 2001,
				},
			},
			NotificationEmails: []string{"test_email1@akamai.com", "test_email2@akamai.com"},
			GroupID:            1337,
			Properties: []datastream.PropertyID{
				{
					PropertyID: 1,
				},
				{
					PropertyID: 2,
				},
				{
					PropertyID: 3,
				},
			},
			StreamName: "test_stream",
		}

		createStreamRequest := datastream.CreateStreamRequest{
			StreamConfiguration: streamConfiguration,
			Activate:            true,
		}

		updateStreamResponse := &datastream.DetailedStreamVersion{
			StreamName:         streamConfiguration.StreamName,
			StreamID:           streamID,
			StreamVersion:      1,
			GroupID:            streamConfiguration.GroupID,
			ContractID:         streamConfiguration.ContractID,
			NotificationEmails: streamConfiguration.NotificationEmails,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
				{
					PropertyID:   2,
					PropertyName: "property_2",
				},
				{
					PropertyID:   3,
					PropertyName: "property_3",
				},
			},
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
				{
					DatasetFieldID:          1002,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
				{
					DatasetFieldID:          2000,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
				{
					DatasetFieldID:          2001,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
			},
			Destination: datastream.Destination{
				Bucket:          "s3_test_bucket",
				DestinationType: datastream.DestinationTypeS3,
				DisplayName:     "s3_test_connector_name",
				Path:            "s3_test_path",
				Region:          "s3_test_region",
			},
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			LatestVersion:         1,
			ProductID:             "Download_Delivery",
			CreatedBy:             "johndoe",
			CreatedDate:           "10-07-2020 12:19:02 GMT",
			StreamStatus:          datastream.StreamStatusActivating,
			ModifiedBy:            "janesmith",
			ModifiedDate:          "15-07-2020 05:51:52 GMT",
		}

		updateStreamRequest := datastream.UpdateStreamRequest{
			StreamID: 12321,
			Activate: true,
			StreamConfiguration: datastream.StreamConfiguration{
				DeliveryConfiguration: datastream.DeliveryConfiguration{
					Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
					Format:    datastream.FormatTypeStructured,
					Frequency: datastream.Frequency{
						IntervalInSeconds: datastream.IntervalInSeconds30,
					},
					UploadFilePrefix: "prefix_updated",
					UploadFileSuffix: "suf_updated",
				},
				Destination: datastream.AbstractConnector(
					&datastream.S3Connector{
						AccessKey:       "s3_test_access_key",
						Bucket:          "s3_test_bucket_updated",
						DisplayName:     "s3_test_connector_name_updated",
						Path:            "s3_test_path",
						Region:          "s3_test_region",
						SecretAccessKey: "s3_test_secret_key",
					},
				),
				ContractID: streamConfiguration.ContractID,
				DatasetFields: []datastream.DatasetFieldID{
					{
						DatasetFieldID: 2000,
					},
					{
						DatasetFieldID: 1002,
					},
					{
						DatasetFieldID: 2001,
					},
					{
						DatasetFieldID: 1001,
					},
				},
				NotificationEmails: []string{"test_email1_updated@akamai.com", "test_email2@akamai.com"},
				Properties:         streamConfiguration.Properties,
				StreamName:         "test_stream_with_updated",
			},
		}

		modifyResponse := func(r datastream.DetailedStreamVersion, opt func(r *datastream.DetailedStreamVersion)) *datastream.DetailedStreamVersion {
			opt(&r)
			return &r
		}

		getStreamResponseActivated := &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusActivated,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination{
				Bucket:          "s3_test_bucket",
				DestinationType: datastream.DestinationTypeS3,
				DisplayName:     "s3_test_connector_name",
				Path:            "s3_test_path",
				Region:          "s3_test_region",
			},
			ContractID:  streamConfiguration.ContractID,
			CreatedBy:   "johndoe",
			CreatedDate: "10-07-2020 12:19:02 GMT",
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
				{
					DatasetFieldID:          1002,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
				{
					DatasetFieldID:          2000,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
				{
					DatasetFieldID:          2001,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
			},
			NotificationEmails: streamConfiguration.NotificationEmails,
			GroupID:            streamConfiguration.GroupID,
			ModifiedBy:         "janesmith",
			ModifiedDate:       "15-07-2020 05:51:52 GMT",
			ProductID:          "Download_Delivery",
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
				{
					PropertyID:   2,
					PropertyName: "property_2",
				},
				{
					PropertyID:   3,
					PropertyName: "property_3",
				},
			},
			StreamID:      updateStreamResponse.StreamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: updateStreamResponse.StreamVersion,
		}

		getStreamResponseStreamActivating := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.StreamStatus = datastream.StreamStatusActivating
		})

		getStreamResponseStreamActivatingAfterUpdate := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.DeliveryConfiguration = datastream.DeliveryConfiguration{
				Delimiter:        updateStreamRequest.StreamConfiguration.DeliveryConfiguration.Delimiter,
				Format:           updateStreamRequest.StreamConfiguration.DeliveryConfiguration.Format,
				Frequency:        updateStreamRequest.StreamConfiguration.DeliveryConfiguration.Frequency,
				UploadFilePrefix: updateStreamRequest.StreamConfiguration.DeliveryConfiguration.UploadFilePrefix,
				UploadFileSuffix: updateStreamRequest.StreamConfiguration.DeliveryConfiguration.UploadFileSuffix,
			}
			r.NotificationEmails = updateStreamRequest.StreamConfiguration.NotificationEmails
			r.StreamName = updateStreamRequest.StreamConfiguration.StreamName
			r.Destination = datastream.Destination{

				Bucket:          "s3_test_bucket_updated",
				DestinationType: datastream.DestinationTypeS3,
				DisplayName:     "s3_test_connector_name_updated",
				Path:            "s3_test_path",
				Region:          "s3_test_region",
			}
			r.DatasetFields = []datastream.DataSetField{
				{
					DatasetFieldID:          2000,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
				{
					DatasetFieldID:          1002,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
				{
					DatasetFieldID:          2001,
					DatasetFieldName:        "dataset_field_name_2",
					DatasetFieldDescription: "dataset_field_desc_2",
				},
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			}
		})

		getStreamResponseStreamActivatedAfterUpdate := modifyResponse(*getStreamResponseStreamActivatingAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.StreamStatus = datastream.StreamStatusActivated
		})

		getStreamResponseDeactivating := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.StreamStatus = datastream.StreamStatusDeactivating
		})

		getStreamResponseDeactivated := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.StreamStatus = datastream.StreamStatusDeactivated
		})

		getStreamRequest := datastream.GetStreamRequest{
			StreamID: streamID,
		}

		client.On("CreateStream", mock.Anything, createStreamRequest).
			Return(updateStreamResponse, nil)

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivating, nil).
			Times(3)

		// first for finishing waitForStreamStatusChange
		// second for complete CreateStream
		// third for reading resource state
		// fourth for reading stream status in UpdateStream
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseActivated, nil).
			Times(4)

		client.On("UpdateStream", mock.Anything, updateStreamRequest).
			Return(updateStreamResponse, nil).
			Once()

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivatedAfterUpdate, nil).
			Times(3)

		// first for finishing waitForStreamStatusChange in UpdateStream
		// second for reading resource state after UpdateStream
		// third for reading stream status in DeleteStream
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivatedAfterUpdate, nil).
			Times(3)

		client.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
			StreamID: streamID,
		}).Return(' ', nil).Once()

		client.On("DeactivateStream", mock.Anything, datastream.DeactivateStreamRequest{
			StreamID: 12321,
		}).Return(&datastream.DetailedStreamVersion{
			StreamID:      streamID,
			StreamVersion: 1,
		}, nil).Once()

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseDeactivating, nil).
			Times(3)

		// for finishing waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseDeactivated, nil).
			Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/lifecycle/create_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "collect_midgress", "false"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.field_delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_prefix", "pre"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_suffix", "suf"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.0.interval_in_secs", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.#", "4"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.0", "1001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.1", "1002"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.2", "2000"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.3", "2001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.0", "test_email1@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "properties.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.display_name", "s3_test_connector_name"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.path", "s3_test_path"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.region", "s3_test_region"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.secret_access_key", "s3_test_secret_key"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/lifecycle/update_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "collect_midgress", "false"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.field_delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_prefix", "prefix_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_suffix", "suf_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.0.interval_in_secs", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.#", "4"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.0", "2000"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.1", "1002"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.2", "2001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.3", "1001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.0", "test_email1_updated@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "properties.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream_with_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.display_name", "s3_test_connector_name_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.path", "s3_test_path"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.region", "s3_test_region"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.secret_access_key", "s3_test_secret_key"),
						),
					},
				},
			})

			client.AssertExpectations(t)
		})
	})
}

func TestResourceUpdate(t *testing.T) {
	PollForActivationStatusChangeInterval = 1 * time.Millisecond
	tests := map[string]struct {
		CreateStreamActive bool
		UpdateStreamActive bool
	}{
		"create active, update active": {
			CreateStreamActive: true,
			UpdateStreamActive: true,
		},
		"create active, update inactive": {
			CreateStreamActive: true,
			UpdateStreamActive: false,
		},
		"create inactive, update active": {
			CreateStreamActive: false,
			UpdateStreamActive: true,
		},
	}
	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
			Format:    datastream.FormatTypeStructured,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
			UploadFilePrefix: "pre",
			UploadFileSuffix: "suf",
		},
		Destination: datastream.AbstractConnector(
			&datastream.OracleCloudStorageConnector{
				AccessKey:       "access_key",
				Bucket:          "bucket",
				DisplayName:     "display_name",
				Namespace:       "namespace",
				Path:            "path",
				Region:          "region",
				SecretAccessKey: "secret_access_key",
			},
		),
		ContractID: "test_contract",
		DatasetFields: []datastream.DatasetFieldID{
			{
				DatasetFieldID: 1001,
			},
		},
		GroupID: 1337,
		Properties: []datastream.PropertyID{
			{
				PropertyID: 1,
			},
		},
		StreamName: "test_stream",
	}
	streamConfigurationFactory := func() datastream.StreamConfiguration {
		return streamConfiguration
	}

	createStreamRequestFactory := func(activateNow bool) datastream.CreateStreamRequest {
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationFactory(),
			Activate:            activateNow,
		}
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamName:    streamConfiguration.StreamName,
		StreamID:      streamID,
		StreamVersion: 2,
		GroupID:       streamConfiguration.GroupID,
		ContractID:    streamConfiguration.ContractID,
		Properties: []datastream.Property{
			{
				PropertyID:   1,
				PropertyName: "property_1",
			},
		},
		DatasetFields: []datastream.DataSetField{
			{
				DatasetFieldID:          1001,
				DatasetFieldName:        "dataset_field_name_1",
				DatasetFieldDescription: "dataset_field_desc_1",
			},
		},
		Destination: datastream.Destination{
			DestinationType: datastream.DestinationTypeOracle,
			Bucket:          "bucket",
			DisplayName:     "display_name",
			Namespace:       "namespace",
			Path:            "path",
			Region:          "region",
		},
		DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
		LatestVersion:         2,
		ProductID:             "Download_Delivery",
		CreatedBy:             "johndoe",
		CreatedDate:           "10-07-2020 12:19:02 GMT",
		StreamStatus:          datastream.StreamStatusActivating,
		ModifiedBy:            "janesmith",
		ModifiedDate:          "15-07-2020 05:51:52 GMT",
	}

	responseFactory := func(activationStatus datastream.StreamStatus) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus: activationStatus,
			DeliveryConfiguration: datastream.DeliveryConfiguration{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Destination: datastream.Destination{
				Bucket:          "bucket",
				DisplayName:     "display_name",
				DestinationType: datastream.DestinationTypeOracle,
				Namespace:       "namespace",
				Path:            "path",
				Region:          "region",
			},
			ContractID: "test_contract",
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			},
			GroupID: 1337,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    "test_stream",
			StreamVersion: 1,
			LatestVersion: 1,
			ProductID:     "API_Acceleration",
		}
	}

	commonChecks := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.field_delimiter", string(datastream.DelimiterTypeSpace)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.format", string(datastream.FormatTypeStructured)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_prefix", "pre"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.upload_file_suffix", "suf"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "delivery_configuration.0.frequency.0.interval_in_secs", "30"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "properties.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.access_key", "access_key"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.bucket", "bucket"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.compress_logs", "false"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.display_name", "display_name"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.namespace", "namespace"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.path", "path"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.region", "region"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.secret_access_key", "secret_access_key"),
	)

	type mockConfig struct {
		status  datastream.StreamStatus
		repeats int
	}

	configureMock := func(m *datastream.Mock, statuses ...mockConfig) {
		for _, statusConfig := range statuses {
			m.On("GetStream", mock.Anything, mock.Anything).
				Return(responseFactory(statusConfig.status), nil).
				Times(statusConfig.repeats)
		}
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &datastream.Mock{}
			m.On("CreateStream", mock.Anything, createStreamRequestFactory(test.CreateStreamActive)).
				Return(updateStreamResponse, nil).
				Once()

			createStreamFilenameSuffix, updateStreamFilenameSuffix := "", ""
			if test.CreateStreamActive {
				createStreamFilenameSuffix = "active"

				configureMock(m, []mockConfig{
					{status: datastream.StreamStatusActivating, repeats: 1},
					{status: datastream.StreamStatusActivated, repeats: 5},
				}...)
			} else {
				createStreamFilenameSuffix = "inactive"

				configureMock(m, []mockConfig{
					{status: datastream.StreamStatusInactive, repeats: 6},
				}...)
			}

			if test.UpdateStreamActive {
				updateStreamFilenameSuffix = "active"

				if test.CreateStreamActive {
					configureMock(m, []mockConfig{
						{status: datastream.StreamStatusActivated, repeats: 2},
						{status: datastream.StreamStatusActivating, repeats: 1},
						{status: datastream.StreamStatusActivated, repeats: 3},
					}...)
				} else {
					m.On("ActivateStream", mock.Anything, mock.Anything).
						Return(&datastream.DetailedStreamVersion{
							StreamVersion: updateStreamResponse.StreamVersion,
						}, nil).
						Once()

					configureMock(m, []mockConfig{
						{status: datastream.StreamStatusActivating, repeats: 1},
						{status: datastream.StreamStatusActivated, repeats: 5},
					}...)
				}
			} else {
				updateStreamFilenameSuffix = "inactive"

				if test.CreateStreamActive {
					configureMock(m, []mockConfig{
						{status: datastream.StreamStatusActivated, repeats: 2},
						{status: datastream.StreamStatusDeactivating, repeats: 1},
						{status: datastream.StreamStatusDeactivated, repeats: 4},
					}...)
				}
			}

			// DeleteStream method will deactivate the stream
			m.On("DeactivateStream", mock.Anything, mock.Anything).
				Return(&datastream.DetailedStreamVersion{
					StreamVersion: updateStreamResponse.StreamVersion,
				}, nil).
				Once()

			// waitForStreamStatusChange in DeleteStream
			configureMock(m, []mockConfig{
				{status: datastream.StreamStatusDeactivating, repeats: 1},
				{status: datastream.StreamStatusDeactivated, repeats: 1},
			}...)

			m.On("DeleteStream", mock.Anything, mock.Anything).
				Return(' ', nil).
				Once()

			useClient(m, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/update_resource/create_stream_%s.tf", createStreamFilenameSuffix)),
							Check: resource.ComposeTestCheckFunc(
								commonChecks,
								resource.TestCheckResourceAttr("akamai_datastream.s", "active", strconv.FormatBool(test.CreateStreamActive)),
							),
						},
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/update_resource/update_stream_%s.tf", updateStreamFilenameSuffix)),
							Check: resource.ComposeTestCheckFunc(
								commonChecks,
								resource.TestCheckResourceAttr("akamai_datastream.s", "active", strconv.FormatBool(test.UpdateStreamActive)),
							),
						},
					},
				})

				m.AssertExpectations(t)
			})
		})
	}
}

func TestResourceStreamErrors(t *testing.T) {
	tests := map[string]struct {
		tfFile    string
		init      func(*datastream.Mock)
		withError *regexp.Regexp
	}{
		"missing required argument": {
			tfFile:    "testdata/TestResourceStream/errors/missing_required_argument/missing_required_argument.tf",
			withError: regexp.MustCompile("Missing required argument"),
		},
		"internal server error": {
			tfFile: "testdata/TestResourceStream/errors/internal_server_error/internal_server_error.tf",
			init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("%w: request failed: %s", datastream.ErrCreateStream, errors.New("500")))
			},
			withError: regexp.MustCompile(datastream.ErrCreateStream.Error()),
		},
		"stream with this name already exists": {
			tfFile: "testdata/TestResourceStream/errors/stream_name_not_unique/stream_name_not_unique.tf",
			init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("%s: %w", datastream.ErrCreateStream, &datastream.Error{
						Type:       "bad-request",
						Title:      "Bad Request",
						StatusCode: 400,
						Errors: []datastream.RequestErrors{
							{
								Type:   "bad-request",
								Title:  "Bad Request",
								Detail: "Stream with name test_stream already exists.",
							},
						},
					}))
			},
			withError: regexp.MustCompile("Stream with name test_stream already exists"),
		},
		"invalid email format": {
			tfFile:    "testdata/TestResourceStream/errors/invalid_email/invalid_email.tf",
			withError: regexp.MustCompile("must be a valid email address"),
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

func TestResourceStreamCustomDiff(t *testing.T) {
	client := &datastream.Mock{}

	tests := map[string]struct {
		tfFile    string
		withError *regexp.Regexp
	}{
		"prefix and suffix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff1.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"prefix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff2.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"suffix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff3.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"prefix and suffix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff4.tf",
			withError: nil,
		},
		"prefix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff5.tf",
			withError: nil,
		},
		"suffix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff6.tf",
			withError: nil,
		},
		"prefix and suffix not present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff7.tf",
			withError: nil,
		},
		"prefix and suffix not present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff8.tf",
			withError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:             testutils.LoadFixtureString(t, test.tfFile),
							ExpectError:        test.withError,
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})
		})

		client.AssertExpectations(t)
	}
}

func TestEmailIDs(t *testing.T) {
	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
			Format:    datastream.FormatTypeStructured,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
			//	UploadFilePrefix: DefaultUploadFilePrefix,
			//	UploadFileSuffix: DefaultUploadFileSuffix,
		},
		Destination: datastream.AbstractConnector(
			&datastream.SplunkConnector{
				CompressLogs:        false,
				DisplayName:         "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				Endpoint:            "splunk_url",
			},
		),
		ContractID: "test_contract",
		DatasetFields: []datastream.DatasetFieldID{
			{
				DatasetFieldID: 1001,
			},
		},
		GroupID: 1337,
		Properties: []datastream.PropertyID{
			{
				PropertyID: 1,
			},
		},
		StreamName: "test_stream",
	}

	createStreamRequestFactory := func(emailIDs []string) datastream.CreateStreamRequest {
		streamConfigurationWithEmailIDs := streamConfiguration
		if emailIDs != nil && len(emailIDs) != 0 {
			streamConfigurationWithEmailIDs.NotificationEmails = emailIDs
		}
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithEmailIDs,
			Activate:            false,
		}
	}

	responseFactory := func(emailIDs []string) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination{
				DestinationType: datastream.DestinationTypeSplunk,
				CompressLogs:    false,
				DisplayName:     "splunk_test_connector_name",
				Endpoint:        "splunk_url",
			},
			ContractID: streamConfiguration.ContractID,
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID: 1001,
				},
			},
			NotificationEmails: emailIDs,
			GroupID:            streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: 2,
		}
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamID:      streamID,
		StreamVersion: 1,
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	tests := map[string]struct {
		Filename   string
		Response   *datastream.DetailedStreamVersion
		EmailIDs   []string
		TestChecks []resource.TestCheckFunc
	}{
		"two emails": {
			Filename: "two_emails.tf",
			EmailIDs: []string{"test_email1@akamai.com", "test_email2@akamai.com"},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "2"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.0", "test_email1@akamai.com"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.1", "test_email2@akamai.com"),
			},
		},
		"one email": {
			Filename: "one_email.tf",
			EmailIDs: []string{"test_email1@akamai.com"},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.0", "test_email1@akamai.com"),
			},
		},
		"empty email": {
			Filename: "empty_email_ids.tf",
			EmailIDs: []string{},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "0"),
			},
		},
		"no email_ids field": {
			Filename: "no_email_ids.tf",
			EmailIDs: []string{},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "notification_emails.#", "0"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}

			createStreamRequest := createStreamRequestFactory(test.EmailIDs)
			client.On("CreateStream", mock.Anything, createStreamRequest).
				Return(updateStreamResponse, nil)

			getStreamResponse := responseFactory(test.EmailIDs)
			client.On("GetStream", mock.Anything, getStreamRequest).
				Return(getStreamResponse, nil)

			client.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
				StreamID: streamID,
			}).Return(' ', nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/email_ids/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}

}

func TestDatasetIDsDiff(t *testing.T) {
	tests := map[string]struct {
		preConfig             string
		fileDatasetIdsOrder   []datastream.DatasetFieldID
		serverDatasetIdsOrder []int
		format                datastream.FormatType
		expectNonEmptyPlan    bool
	}{
		"no order mixed in json config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1002, 1001},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    false,
		},
		"id change in json config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1002, 1003},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"duplicates in server side json config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1002, 1002},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"duplicates in incoming json config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/json_config_duplicates.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1002,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1001, 1002},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"no order mixed in structured config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/structured_config.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1002, 1001},
			format:                datastream.FormatTypeStructured,
			expectNonEmptyPlan:    true,
		},
		"id change in structured config": {
			preConfig: "testdata/TestResourceStream/dataset_ids_diff/structured_config.tf",
			fileDatasetIdsOrder: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				}, {
					DatasetFieldID: 1002,
				},
			},
			serverDatasetIdsOrder: []int{1002, 1003},
			format:                datastream.FormatTypeStructured,
			expectNonEmptyPlan:    true,
		},
	}

	for name, test := range tests {

		streamConfiguration := datastream.StreamConfiguration{

			DeliveryConfiguration: datastream.DeliveryConfiguration{
				Format: test.format,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
				},
			},
			Destination: datastream.AbstractConnector(
				&datastream.SplunkConnector{
					CompressLogs:        false,
					DisplayName:         "splunk_test_connector_name",
					EventCollectorToken: "splunk_event_collector_token",
					Endpoint:            "splunk_url",
				},
			),
			ContractID:    "test_contract",
			DatasetFields: test.fileDatasetIdsOrder,
			GroupID:       1337,
			Properties: []datastream.PropertyID{
				{
					PropertyID: 1,
				},
			},
			StreamName: "test_stream",
		}

		if test.format == datastream.FormatTypeStructured {
			streamConfiguration.DeliveryConfiguration.Delimiter = datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace)
		}

		createStreamRequest := datastream.CreateStreamRequest{
			StreamConfiguration: streamConfiguration,
			Activate:            false,
		}

		createStreamResponse := &datastream.DetailedStreamVersion{
			StreamID:      streamID,
			StreamVersion: 1,
		}

		getStreamRequest := datastream.GetStreamRequest{
			StreamID: streamID,
		}

		getStreamResponse := &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination{
				DestinationType: datastream.DestinationTypeSplunk,
				CompressLogs:    false,
				DisplayName:     "splunk_test_connector_name",
				Endpoint:        "splunk_url",
			},
			ContractID: streamConfiguration.ContractID,
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID: test.serverDatasetIdsOrder[0],
				},
				{
					DatasetFieldID: test.serverDatasetIdsOrder[1],
				},
			},
			GroupID: streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: 2,
		}

		deleteStreamRequest := datastream.DeleteStreamRequest{
			StreamID: streamID,
		}

		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}

			client.On("CreateStream", mock.Anything, createStreamRequest).
				Return(createStreamResponse, nil)

			client.On("GetStream", mock.Anything, getStreamRequest).
				Return(getStreamResponse, nil).Times(3)

			client.On("DeleteStream", mock.Anything, deleteStreamRequest).
				Return(' ', nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:             testutils.LoadFixtureString(t, test.preConfig),
							ExpectNonEmptyPlan: test.expectNonEmptyPlan,
							Check: resource.ComposeTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields.#", "2"),
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields.0", strconv.Itoa(test.serverDatasetIdsOrder[0])),
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields.1", strconv.Itoa(test.serverDatasetIdsOrder[1])),
							),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})

	}
}

func TestCustomHeaders(t *testing.T) {
	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
			//UploadFilePrefix: DefaultUploadFilePrefix,
			//UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID: "test_contract",
		DatasetFields: []datastream.DatasetFieldID{
			{
				DatasetFieldID: 1001,
			},
		},
		GroupID: 1337,
		Properties: []datastream.PropertyID{
			{
				PropertyID: 1,
			},
		},
		StreamName: "test_stream",
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Destination = datastream.AbstractConnector(
			connector,
		)
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
			Activate:            false,
		}
	}

	responseFactory := func(connector datastream.Destination) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination(
				connector,
			),
			ContractID: streamConfiguration.ContractID,
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			},
			GroupID: streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: 2,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamID:      streamID,
		StreamVersion: 1,
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.Destination
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"splunk": {
			Filename: "custom_headers_splunk.tf",
			Connector: &datastream.SplunkConnector{
				CompressLogs:        false,
				DisplayName:         "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				Endpoint:            "splunk_url",
				CustomHeaderName:    "custom_header_name",
				CustomHeaderValue:   "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:   datastream.DestinationTypeSplunk,
				CompressLogs:      false,
				DisplayName:       "splunk_test_connector_name",
				Endpoint:          "splunk_url",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.display_name", "splunk_test_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.event_collector_token", "splunk_event_collector_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.endpoint", "splunk_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.custom_header_value", "custom_header_value"),
			},
		},
		"https": {
			Filename: "custom_headers_https.tf",
			Connector: &datastream.CustomHTTPSConnector{
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				DisplayName:        "HTTPS connector name",
				Password:           "password",
				Endpoint:           "https_connector_url",
				UserName:           "username",
				ContentType:        "content_type",
				CustomHeaderName:   "custom_header_name",
				CustomHeaderValue:  "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:    datastream.DestinationTypeHTTPS,
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				DisplayName:        "HTTPS connector name",
				Endpoint:           "https_connector_url",
				ContentType:        "content_type",
				CustomHeaderName:   "custom_header_name",
				CustomHeaderValue:  "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.authentication_type", "BASIC"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.display_name", "HTTPS connector name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.compress_logs", "true"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.custom_header_value", "custom_header_value"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.endpoint", "https_connector_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.user_name", "username"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.password", "password"),
			},
		},
		"sumologic": {
			Filename: "custom_headers_sumologic.tf",
			Connector: &datastream.SumoLogicConnector{
				CollectorCode:     "collector_code",
				CompressLogs:      true,
				DisplayName:       "Sumologic connector name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:   datastream.DestinationTypeSumoLogic,
				CompressLogs:      true,
				DisplayName:       "Sumologic connector name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "Sumologic connector name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.compress_logs", "true"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_value", "custom_header_value"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
			},
		},
		"loggly": {
			Filename: "custom_headers_loggly.tf",
			Connector: &datastream.LogglyConnector{
				DisplayName:       "loggly_connector_name",
				Endpoint:          "endpoint",
				AuthToken:         "auth_token",
				Tags:              "tag1,tag2,tag3",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:   datastream.DestinationTypeLoggly,
				DisplayName:       "loggly_connector_name",
				Endpoint:          "endpoint",
				Tags:              "tag1,tag2,tag3",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.display_name", "loggly_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.endpoint", "endpoint"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.auth_token", "auth_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.tags", "tag1,tag2,tag3"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.custom_header_value", "custom_header_value"),
			},
		},
		"new_relic": {
			Filename: "custom_headers_new_relic.tf",
			Connector: &datastream.NewRelicConnector{
				DisplayName:       "new_relic_connector_name",
				Endpoint:          "endpoint",
				AuthToken:         "auth_token",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:   datastream.DestinationTypeNewRelic,
				DisplayName:       "new_relic_connector_name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.display_name", "new_relic_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.endpoint", "endpoint"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.auth_token", "auth_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.custom_header_value", "custom_header_value"),
			},
		},
		"elasticsearch": {
			Filename: "custom_headers_elasticsearch.tf",
			Connector: &datastream.ElasticsearchConnector{
				DisplayName:       "elasticsearch_connector_name",
				Endpoint:          "endpoint",
				IndexName:         "index_name",
				UserName:          "user_name",
				Password:          "password",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.Destination{
				DestinationType:   datastream.DestinationTypeElasticsearch,
				DisplayName:       "elasticsearch_connector_name",
				Endpoint:          "endpoint",
				IndexName:         "index_name",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.display_name", "elasticsearch_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.endpoint", "endpoint"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.custom_header_value", "custom_header_value"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}

			createStreamRequest := createStreamRequestFactory(test.Connector)
			client.On("CreateStream", mock.Anything, createStreamRequest).
				Return(updateStreamResponse, nil)

			getStreamResponse := responseFactory(test.Response)
			client.On("GetStream", mock.Anything, getStreamRequest).
				Return(getStreamResponse, nil)

			client.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
				StreamID: streamID,
			}).Return(' ', nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/custom_headers/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}
}

func TestMTLS(t *testing.T) {
	streamID := int64(12321)

	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
			//UploadFilePrefix: DefaultUploadFilePrefix,
			//UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID: "test_contract",
		DatasetFields: []datastream.DatasetFieldID{
			{
				DatasetFieldID: 1001,
			},
		},
		GroupID: 1337,
		Properties: []datastream.PropertyID{
			{
				PropertyID: 1,
			},
		},
		StreamName: "test_stream",
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Destination = datastream.AbstractConnector(
			connector,
		)
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
			Activate:            false,
		}
	}

	responseFactory := func(connector datastream.Destination) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination(
				connector,
			),
			ContractID: streamConfiguration.ContractID,
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			},
			GroupID: streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: 2,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamID:      streamID,
		StreamVersion: 1,
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.Destination
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"splunk_mtls": {
			Filename: "mtls_splunk.tf",
			Connector: &datastream.SplunkConnector{
				CompressLogs:        false,
				DisplayName:         "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				Endpoint:            "splunk_url",
				TLSHostname:         "tls_hostname",
				CACert:              "ca_cert",
				ClientCert:          "client_cert",
				ClientKey:           "client_key",
			},
			Response: datastream.Destination{
				DestinationType: datastream.DestinationTypeSplunk,
				CompressLogs:    false,
				DisplayName:     "splunk_test_connector_name",
				Endpoint:        "splunk_url",
				TLSHostname:     "tls_hostname",
				MTLS:            "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.display_name", "splunk_test_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.event_collector_token", "splunk_event_collector_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.endpoint", "splunk_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.tls_hostname", "tls_hostname"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.ca_cert", "ca_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.client_cert", "client_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.client_key", "client_key"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.m_tls", "true"),
			},
		},
		"https_mtls": {
			Filename: "mtls_https.tf",
			Connector: &datastream.CustomHTTPSConnector{
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				DisplayName:        "HTTPS connector name",
				Password:           "password",
				Endpoint:           "https_connector_url",
				UserName:           "username",
				ContentType:        "content_type",
				TLSHostname:        "tls_hostname",
				CACert:             "ca_cert",
				ClientCert:         "client_cert",
				ClientKey:          "client_key",
			},
			Response: datastream.Destination{
				DestinationType:    datastream.DestinationTypeHTTPS,
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				DisplayName:        "HTTPS connector name",
				Endpoint:           "https_connector_url",
				ContentType:        "content_type",
				TLSHostname:        "tls_hostname",
				MTLS:               "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.authentication_type", "BASIC"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.display_name", "HTTPS connector name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.compress_logs", "true"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.endpoint", "https_connector_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.user_name", "username"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.password", "password"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.tls_hostname", "tls_hostname"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.ca_cert", "ca_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.client_cert", "client_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.client_key", "client_key"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.m_tls", "true"),
			},
		},
		"elasticsearch_mtls": {
			Filename: "mtls_elasticsearch.tf",
			Connector: &datastream.ElasticsearchConnector{
				DisplayName: "elasticsearch_connector_name",
				Endpoint:    "endpoint",
				IndexName:   "index_name",
				UserName:    "user_name",
				Password:    "password",
				ContentType: "content_type",
				TLSHostname: "tls_hostname",
				CACert:      "ca_cert",
				ClientCert:  "client_cert",
				ClientKey:   "client_key",
			},
			Response: datastream.Destination{
				DestinationType: datastream.DestinationTypeElasticsearch,
				CompressLogs:    true,
				DisplayName:     "elasticsearch_connector_name",
				Endpoint:        "endpoint",
				IndexName:       "index_name",
				ContentType:     "content_type",
				TLSHostname:     "tls_hostname",
				MTLS:            "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.display_name", "elasticsearch_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.endpoint", "endpoint"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.index_name", "index_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.user_name", "user_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.password", "password"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.tls_hostname", "tls_hostname"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.ca_cert", "ca_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.client_cert", "client_cert"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.client_key", "client_key"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.m_tls", "true"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}

			createStreamRequest := createStreamRequestFactory(test.Connector)
			client.On("CreateStream", mock.Anything, createStreamRequest).
				Return(updateStreamResponse, nil)

			getStreamResponse := responseFactory(test.Response)
			client.On("GetStream", mock.Anything, getStreamRequest).
				Return(getStreamResponse, nil)

			client.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
				StreamID: streamID,
			}).Return(' ', nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/mtls/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}
}

/* commenting this as it is not required in V2
func TestUrlSuppressor(t *testing.T) {

	streamConfigurationFactory := func(connector datastream.AbstractConnector) datastream.StreamConfiguration {
		return datastream.StreamConfiguration{
			DeliveryConfiguration: datastream.DeliveryConfiguration{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
				},
				UploadFilePrefix: "ak",
				UploadFileSuffix: "ds",
			},
			Destination: datastream.AbstractConnector(
				connector,
			),
			ContractID: "test_contract",
			DatasetFields: []datastream.DatasetFieldID{
				{
					DatasetFieldID: 1001,
				},
			},
			GroupID: 1337,
			Properties: []datastream.PropertyID{
				{
					PropertyID: 1,
				},
			},
			StreamName: "test_stream",
		}
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationFactory(connector),
			Activate:            false,
		}
	}

	updateStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.UpdateStreamRequest {
		req := datastream.UpdateStreamRequest{
			StreamID:            streamID,
			StreamConfiguration: streamConfigurationFactory(connector),
		}
		req.StreamConfiguration.GroupID = 1337
		return req
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamID:      streamID,
		StreamVersion: 1,
	}

	responseFactory := func(connector datastream.Destination) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus: datastream.StreamStatusInactive,
			DeliveryConfiguration: datastream.DeliveryConfiguration{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
				},
			},
			Destination: connector,
			ContractID:  "test_contract",
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			},
			GroupID: 1337,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    "test_stream",
			StreamVersion: 1,
		}
	}

	tests := map[string]struct {
		Init  func(m *datastream.Mock)
		Steps []resource.TestStep
	}{
		"idempotent when endpoint is stripped by api": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					DisplayName:   "display_name",
					Endpoint:      "endpoint/?/?",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeSumoLogic,
						CompressLogs:    true,
						DisplayName:     "display_name",
						Endpoint:        "endpoint", //api returns stripped url
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/idempotency/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/idempotency/create_stream.tf"),
					PlanOnly: true,
				},
			},
		},
		"update endpoint field": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					DisplayName:   "display_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeSumoLogic,
						CompressLogs:    true,
						DisplayName:     "display_name",
						Endpoint:        "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					DisplayName:   "display_name",
					Endpoint:      "endpoint_updated",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeSumoLogic,
						CompressLogs:    true,
						DisplayName:     "display_name",
						Endpoint:        "endpoint_updated",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/update_endpoint_field/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/update_endpoint_field/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint_updated"),
					),
				},
			},
		},
		"adding new fields": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					DisplayName:   "display_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeSumoLogic,
						CompressLogs:    true,
						DisplayName:     "display_name",
						Endpoint:        "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode:     "collector_code",
					CompressLogs:      true,
					DisplayName:       "display_name",
					Endpoint:          "endpoint",
					CustomHeaderName:  "custom_header_name",
					CustomHeaderValue: "custom_header_value",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType:   datastream.DestinationTypeSumoLogic,
						CompressLogs:      true,
						DisplayName:       "display_name",
						Endpoint:          "endpoint",
						CustomHeaderName:  "custom_header_name",
						CustomHeaderValue: "custom_header_value",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/adding_fields/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/adding_fields/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_name", "custom_header_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_value", "custom_header_value"),
					),
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/adding_fields/update_stream.tf"),
					PlanOnly: true,
				},
			},
		},
		"change connector": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					DisplayName:   "display_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeSumoLogic,
						CompressLogs:    true,
						DisplayName:     "display_name",
						Endpoint:        "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.DatadogConnector{
					AuthToken:   "auth_token",
					DisplayName: "display_name",
					Endpoint:    "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.Destination{
						DestinationType: datastream.DestinationTypeDataDog,
						DisplayName:     "display_name",
						Endpoint:        "endpoint",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/change_connector/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/change_connector/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.display_name", "display_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.auth_token", "auth_token"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceStream/urlSuppressor/change_connector/update_stream.tf"),
					PlanOnly: true,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &datastream.Mock{}
			test.Init(m)

			m.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
				StreamID: streamID,
			}).Return(' ', nil)

			useClient(m, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:             test.Steps,
				})

				m.AssertExpectations(t)
			})
		})
	}
}*/

func TestConnectors(t *testing.T) {
	streamConfiguration := datastream.StreamConfiguration{
		DeliveryConfiguration: datastream.DeliveryConfiguration{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds30,
			},
			UploadFilePrefix: DefaultUploadFilePrefix,
			UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID: "test_contract",
		DatasetFields: []datastream.DatasetFieldID{
			{
				DatasetFieldID: 1001,
			},
		},
		GroupID: 1337,
		Properties: []datastream.PropertyID{
			{
				PropertyID: 1,
			},
		},
		StreamName: "test_stream",
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Destination = datastream.AbstractConnector(
			connector,
		)
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
			Activate:            false,
		}
	}

	responseFactory := func(connector datastream.Destination) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			StreamStatus:          datastream.StreamStatusInactive,
			DeliveryConfiguration: streamConfiguration.DeliveryConfiguration,
			Destination: datastream.Destination(
				connector,
			),
			ContractID: streamConfiguration.ContractID,
			DatasetFields: []datastream.DataSetField{
				{
					DatasetFieldID:          1001,
					DatasetFieldName:        "dataset_field_name_1",
					DatasetFieldDescription: "dataset_field_desc_1",
				},
			},
			GroupID: streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:      streamID,
			StreamName:    streamConfiguration.StreamName,
			StreamVersion: 2,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.DetailedStreamVersion{
		StreamID:      streamID,
		StreamVersion: 1,
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.Destination
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"azure": {
			Filename: "azure.tf",
			Connector: &datastream.AzureConnector{
				AccessKey:     "access_key",
				AccountName:   "account_name",
				DisplayName:   "connector_name",
				ContainerName: "container_name",
				Path:          "path",
			},
			Response: datastream.Destination{
				DestinationType: datastream.DestinationTypeAzure,
				AccountName:     "account_name",
				DisplayName:     "connector_name",
				ContainerName:   "container_name",
				Path:            "path",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.account_name", "account_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.display_name", "connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.container_name", "container_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.path", "path"),
			},
		},
		"gcs": {
			Filename: "gcs.tf",
			Connector: &datastream.GCSConnector{
				Bucket:             "bucket",
				DisplayName:        "connector_name",
				Path:               "path",
				PrivateKey:         "private_key",
				ProjectID:          "project_id",
				ServiceAccountName: "service_account_name",
			},
			Response: datastream.Destination{
				DestinationType:    datastream.DestinationTypeGcs,
				Bucket:             "bucket",
				DisplayName:        "connector_name",
				Path:               "path",
				ProjectID:          "project_id",
				ServiceAccountName: "service_account_name",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.bucket", "bucket"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.display_name", "connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.path", "path"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.project_id", "project_id"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.service_account_name", "service_account_name"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}

			client.On("CreateStream", mock.Anything, createStreamRequestFactory(test.Connector)).
				Return(updateStreamResponse, nil)

			client.On("GetStream", mock.Anything, getStreamRequest).
				Return(responseFactory(test.Response), nil)

			client.On("DeleteStream", mock.Anything, mock.Anything).Return(' ', nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestResourceStream/connectors/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}
}

func TestEmptyFilePrefixSuffixSetForHttpsDestination(t *testing.T) {

	configurationOfPrefixSuffixNotSupportedDest := datastream.DeliveryConfiguration{
		Format: datastream.FormatTypeJson,
		Frequency: datastream.Frequency{
			IntervalInSeconds: datastream.IntervalInSeconds30,
		},
		UploadFilePrefix: DefaultUploadFilePrefix,
		UploadFileSuffix: DefaultUploadFileSuffix,
	}

	result, err := FilePrefixSuffixSet(`splunk_connector`, &configurationOfPrefixSuffixNotSupportedDest)
	require.NoError(t, err)
	assert.Equal(t, "", result.UploadFilePrefix)
	assert.Equal(t, "", result.UploadFileSuffix)
}

func TestFilePrefixSuffixSetForObjectStorageDestination(t *testing.T) {

	configurationOfPrefixSuffixSupportedDest := datastream.DeliveryConfiguration{
		Format: datastream.FormatTypeJson,
		Frequency: datastream.Frequency{
			IntervalInSeconds: datastream.IntervalInSeconds30,
		},
		UploadFilePrefix: "pre",
		UploadFileSuffix: "suf",
	}

	result, err := FilePrefixSuffixSet(`azure_connector`, &configurationOfPrefixSuffixSupportedDest)
	require.NoError(t, err)
	assert.Equal(t, configurationOfPrefixSuffixSupportedDest.UploadFilePrefix, result.UploadFilePrefix)
	assert.Equal(t, configurationOfPrefixSuffixSupportedDest.UploadFileSuffix, result.UploadFileSuffix)
}
