package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

const (
	streamID = int64(12321)
)

func TestResourceStream(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		client := &datastream.Mock{}

		PollForActivationStatusChangeInterval = 1 * time.Millisecond

		streamConfiguration := datastream.StreamConfiguration{
			ActivateNow: true,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Connectors: []datastream.AbstractConnector{
				&datastream.S3Connector{
					AccessKey:       "s3_test_access_key",
					Bucket:          "s3_test_bucket",
					ConnectorName:   "s3_test_connector_name",
					Path:            "s3_test_path",
					Region:          "s3_test_region",
					SecretAccessKey: "s3_test_secret_key",
				},
			},
			ContractID:      "test_contract",
			DatasetFieldIDs: []int{1001, 1002, 2000, 2001},
			EmailIDs:        "test_email1@akamai.com,test_email2@akamai.com",
			GroupID:         tools.IntPtr(1337),
			PropertyIDs:     []int{1, 2, 3},
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}

		createStreamRequest := datastream.CreateStreamRequest{
			StreamConfiguration: streamConfiguration,
		}

		updateStreamResponse := &datastream.StreamUpdate{
			StreamVersionKey: datastream.StreamVersionKey{
				StreamID:        streamID,
				StreamVersionID: 1,
			},
		}

		updateStreamRequest := datastream.UpdateStreamRequest{
			StreamID: 12321,
			StreamConfiguration: datastream.StreamConfiguration{
				ActivateNow: false,
				Config: datastream.Config{
					Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
					Format:    datastream.FormatTypeStructured,
					Frequency: datastream.Frequency{
						TimeInSec: datastream.TimeInSec30,
					},
					UploadFilePrefix: "prefix_updated",
					UploadFileSuffix: "suf_updated",
				},
				Connectors: []datastream.AbstractConnector{
					&datastream.S3Connector{
						AccessKey:       "s3_test_access_key",
						Bucket:          "s3_test_bucket_updated",
						ConnectorName:   "s3_test_connector_name_updated",
						Path:            "s3_test_path",
						Region:          "s3_test_region",
						SecretAccessKey: "s3_test_secret_key",
					},
				},
				ContractID:      streamConfiguration.ContractID,
				DatasetFieldIDs: []int{2000, 1002, 2001, 1001},
				EmailIDs:        "test_email1_updated@akamai.com,test_email2@akamai.com",
				PropertyIDs:     streamConfiguration.PropertyIDs,
				StreamName:      "test_stream_with_updated",
				StreamType:      streamConfiguration.StreamType,
				TemplateName:    streamConfiguration.TemplateName,
			},
		}

		modifyResponse := func(r datastream.DetailedStreamVersion, opt func(r *datastream.DetailedStreamVersion)) *datastream.DetailedStreamVersion {
			opt(&r)
			return &r
		}

		getStreamResponseActivated := &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusActivated,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				{
					Bucket:        "s3_test_bucket",
					ConnectorType: datastream.ConnectorTypeS3,
					ConnectorName: "s3_test_connector_name",
					Path:          "s3_test_path",
					Region:        "s3_test_region",
				},
			},
			ContractID:  streamConfiguration.ContractID,
			CreatedBy:   "johndoe",
			CreatedDate: "10-07-2020 12:19:02 GMT",
			Datasets: []datastream.DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          1001,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   0,
						},
						{
							DatasetFieldID:          1002,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
							Order:                   1,
						},
					},
				},
				{
					DatasetGroupName:        "group_name_2",
					DatasetGroupDescription: "group_desc_2",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          2000,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   2,
						},
						{
							DatasetFieldID:          2001,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
							Order:                   3,
						},
					},
				},
			},
			EmailIDs:     streamConfiguration.EmailIDs,
			Errors:       nil,
			GroupID:      *streamConfiguration.GroupID,
			GroupName:    "Default Group-1-ABCDE",
			ModifiedBy:   "janesmith",
			ModifiedDate: "15-07-2020 05:51:52 GMT",
			ProductID:    "Download_Delivery",
			ProductName:  "Download Delivery",
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
			StreamID:        updateStreamResponse.StreamVersionKey.StreamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: updateStreamResponse.StreamVersionKey.StreamVersionID,
			TemplateName:    streamConfiguration.TemplateName,
		}

		getStreamResponseStreamActivating := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusActivating
		})

		getStreamResponseStreamActivatingAfterUpdate := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.Config = datastream.Config{
				Delimiter:        updateStreamRequest.StreamConfiguration.Config.Delimiter,
				Format:           updateStreamRequest.StreamConfiguration.Config.Format,
				Frequency:        updateStreamRequest.StreamConfiguration.Config.Frequency,
				UploadFilePrefix: updateStreamRequest.StreamConfiguration.Config.UploadFilePrefix,
				UploadFileSuffix: updateStreamRequest.StreamConfiguration.Config.UploadFileSuffix,
			}
			r.EmailIDs = updateStreamRequest.StreamConfiguration.EmailIDs
			r.StreamName = updateStreamRequest.StreamConfiguration.StreamName
			r.Connectors = []datastream.ConnectorDetails{
				{
					Bucket:        "s3_test_bucket_updated",
					ConnectorType: datastream.ConnectorTypeS3,
					ConnectorName: "s3_test_connector_name_updated",
					Path:          "s3_test_path",
					Region:        "s3_test_region",
				},
			}
			r.Datasets = []datastream.DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          1001,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   3,
						},
						{
							DatasetFieldID:          1002,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
							Order:                   1,
						},
					},
				},
				{
					DatasetGroupName:        "group_name_2",
					DatasetGroupDescription: "group_desc_2",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          2000,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   0,
						},
						{
							DatasetFieldID:          2001,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
							Order:                   2,
						},
					},
				},
			}
		})

		getStreamResponseStreamActivatedAfterUpdate := modifyResponse(*getStreamResponseStreamActivatingAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusActivated
		})

		getStreamResponseDeactivating := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusDeactivating
		})

		getStreamResponseDeactivated := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusDeactivated
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
		}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil).Once()

		client.On("DeactivateStream", mock.Anything, datastream.DeactivateStreamRequest{
			StreamID: 12321,
		}).Return(&datastream.DeactivateStreamResponse{
			StreamVersionKey: datastream.StreamVersionKey{
				StreamID:        streamID,
				StreamVersionID: 1,
			},
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResourceStream/lifecycle/create_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_prefix", "pre"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_suffix", "suf"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.0.time_in_sec", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.#", "4"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.0", "1001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.1", "1002"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.2", "2000"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.3", "2001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "property_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_type", string(datastream.StreamTypeRawLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "template_name", string(datastream.TemplateNameEdgeLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.connector_name", "s3_test_connector_name"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.path", "s3_test_path"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.region", "s3_test_region"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.secret_access_key", "s3_test_secret_key"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResourceStream/lifecycle/update_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_prefix", "prefix_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_suffix", "suf_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.0.time_in_sec", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.#", "4"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.0", "2000"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.1", "1002"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.2", "2001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.3", "1001"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1_updated@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "property_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream_with_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_type", string(datastream.StreamTypeRawLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "template_name", string(datastream.TemplateNameEdgeLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.connector_name", "s3_test_connector_name_updated"),
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

	streamConfigurationFactory := func(activateNow bool) datastream.StreamConfiguration {
		return datastream.StreamConfiguration{
			ActivateNow: activateNow,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Connectors: []datastream.AbstractConnector{
				&datastream.OracleCloudStorageConnector{
					AccessKey:       "access_key",
					Bucket:          "bucket",
					ConnectorName:   "connector_name",
					Namespace:       "namespace",
					Path:            "path",
					Region:          "region",
					SecretAccessKey: "secret_access_key",
				},
			},
			ContractID:      "test_contract",
			DatasetFieldIDs: []int{1001},
			GroupID:         tools.IntPtr(1337),
			PropertyIDs:     []int{1},
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}
	}

	createStreamRequestFactory := func(activateNow bool) datastream.CreateStreamRequest {
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationFactory(activateNow),
		}
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 2,
		},
	}

	responseFactory := func(activationStatus datastream.ActivationStatus) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: activationStatus,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Connectors: []datastream.ConnectorDetails{
				{
					Bucket:        "bucket",
					ConnectorName: "connector_name",
					ConnectorType: datastream.ConnectorTypeOracle,
					Namespace:     "namespace",
					Path:          "path",
					Region:        "region",
				},
			},
			ContractID: "test_contract",
			Datasets: []datastream.DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          1001,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   0,
						},
					},
				},
			},
			GroupID: 1337,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			StreamVersionID: 1,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}
	}

	commonChecks := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("akamai_datastream.s", "id", strconv.FormatInt(streamID, 10)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.delimiter", string(datastream.DelimiterTypeSpace)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.format", string(datastream.FormatTypeStructured)),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_prefix", "pre"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_suffix", "suf"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.0.time_in_sec", "30"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "property_ids.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.#", "1"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.access_key", "access_key"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.bucket", "bucket"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.compress_logs", "false"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.connector_name", "connector_name"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.namespace", "namespace"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.path", "path"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.region", "region"),
		resource.TestCheckResourceAttr("akamai_datastream.s", "oracle_connector.0.secret_access_key", "secret_access_key"),
	)

	type mockConfig struct {
		status  datastream.ActivationStatus
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
					{status: datastream.ActivationStatusActivating, repeats: 1},
					{status: datastream.ActivationStatusActivated, repeats: 5},
				}...)
			} else {
				createStreamFilenameSuffix = "inactive"

				configureMock(m, []mockConfig{
					{status: datastream.ActivationStatusInactive, repeats: 6},
				}...)
			}

			if test.UpdateStreamActive {
				updateStreamFilenameSuffix = "active"

				if test.CreateStreamActive {
					configureMock(m, []mockConfig{
						{status: datastream.ActivationStatusActivated, repeats: 2},
						{status: datastream.ActivationStatusActivating, repeats: 1},
						{status: datastream.ActivationStatusActivated, repeats: 3},
					}...)
				} else {
					m.On("ActivateStream", mock.Anything, mock.Anything).
						Return(&datastream.ActivateStreamResponse{
							StreamVersionKey: updateStreamResponse.StreamVersionKey,
						}, nil).
						Once()

					configureMock(m, []mockConfig{
						{status: datastream.ActivationStatusActivating, repeats: 1},
						{status: datastream.ActivationStatusActivated, repeats: 5},
					}...)
				}
			} else {
				updateStreamFilenameSuffix = "inactive"

				if test.CreateStreamActive {
					configureMock(m, []mockConfig{
						{status: datastream.ActivationStatusActivated, repeats: 2},
						{status: datastream.ActivationStatusDeactivating, repeats: 1},
						{status: datastream.ActivationStatusDeactivated, repeats: 4},
					}...)
				}
			}

			// DeleteStream method will deactivate the stream
			m.On("DeactivateStream", mock.Anything, mock.Anything).
				Return(&datastream.DeactivateStreamResponse{
					StreamVersionKey: updateStreamResponse.StreamVersionKey,
				}, nil).
				Once()

			// waitForStreamStatusChange in DeleteStream
			configureMock(m, []mockConfig{
				{status: datastream.ActivationStatusDeactivating, repeats: 1},
				{status: datastream.ActivationStatusDeactivated, repeats: 1},
			}...)

			m.On("DeleteStream", mock.Anything, mock.Anything).
				Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil).
				Once()

			useClient(m, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/update_resource/create_stream_%s.tf", createStreamFilenameSuffix)),
							Check: resource.ComposeTestCheckFunc(
								commonChecks,
								resource.TestCheckResourceAttr("akamai_datastream.s", "active", strconv.FormatBool(test.CreateStreamActive)),
							),
						},
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/update_resource/update_stream_%s.tf", updateStreamFilenameSuffix)),
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

func TestEmailIDs(t *testing.T) {
	streamConfiguration := datastream.StreamConfiguration{
		ActivateNow: false,
		Config: datastream.Config{
			Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
			Format:    datastream.FormatTypeStructured,
			Frequency: datastream.Frequency{
				TimeInSec: datastream.TimeInSec30,
			},
			UploadFilePrefix: DefaultUploadFilePrefix,
			UploadFileSuffix: DefaultUploadFileSuffix,
		},
		Connectors: []datastream.AbstractConnector{
			&datastream.SplunkConnector{
				CompressLogs:        false,
				ConnectorName:       "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				URL:                 "splunk_url",
			},
		},
		ContractID:      "test_contract",
		DatasetFieldIDs: []int{1001},
		GroupID:         tools.IntPtr(1337),
		PropertyIDs:     []int{1},
		StreamName:      "test_stream",
		StreamType:      datastream.StreamTypeRawLogs,
		TemplateName:    datastream.TemplateNameEdgeLogs,
	}

	createStreamRequestFactory := func(emailIDs string) datastream.CreateStreamRequest {
		streamConfigurationWithEmailIDs := streamConfiguration
		if emailIDs != "" {
			streamConfigurationWithEmailIDs.EmailIDs = emailIDs
		}
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithEmailIDs,
		}
	}

	responseFactory := func(emailIDs string) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				{
					ConnectorType: datastream.ConnectorTypeSplunk,
					CompressLogs:  false,
					ConnectorName: "splunk_test_connector_name",
					URL:           "splunk_url",
				},
			},
			ContractID: streamConfiguration.ContractID,
			Datasets: []datastream.DataSets{
				{
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID: 1001,
							Order:          0,
						},
					},
				},
			},
			EmailIDs: emailIDs,
			GroupID:  *streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: 2,
			TemplateName:    streamConfiguration.TemplateName,
		}
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 1,
		},
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	tests := map[string]struct {
		Filename   string
		Response   *datastream.DetailedStreamVersion
		EmailIDs   string
		TestChecks []resource.TestCheckFunc
	}{
		"two emails": {
			Filename: "two_emails.tf",
			EmailIDs: "test_email1@akamai.com,test_email2@akamai.com",
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "2"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1@akamai.com"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.1", "test_email2@akamai.com"),
			},
		},
		"one email": {
			Filename: "one_email.tf",
			EmailIDs: "test_email1@akamai.com",
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1@akamai.com"),
			},
		},
		"empty email": {
			Filename: "empty_email_ids.tf",
			EmailIDs: "",
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "0"),
			},
		},
		"no email_ids field": {
			Filename: "no_email_ids.tf",
			EmailIDs: "",
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "0"),
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
			}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/email_ids/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
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
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.tfFile),
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
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString(test.tfFile),
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

func TestDatasetIDsDiff(t *testing.T) {
	tests := map[string]struct {
		preConfig             string
		fileDatasetIdsOrder   []int
		serverDatasetIdsOrder []int
		format                datastream.FormatType
		expectNonEmptyPlan    bool
	}{
		"order mixed in json config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder:   []int{1001, 1002},
			serverDatasetIdsOrder: []int{1002, 1001},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    false,
		},
		"id change in json config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder:   []int{1001, 1002},
			serverDatasetIdsOrder: []int{1002, 1003},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"duplicates in server side json config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/json_config.tf",
			fileDatasetIdsOrder:   []int{1001, 1002},
			serverDatasetIdsOrder: []int{1002, 1002},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"duplicates in incoming json config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/json_config_duplicates.tf",
			fileDatasetIdsOrder:   []int{1002, 1002},
			serverDatasetIdsOrder: []int{1001, 1002},
			format:                datastream.FormatTypeJson,
			expectNonEmptyPlan:    true,
		},
		"order mixed in structured config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/structured_config.tf",
			fileDatasetIdsOrder:   []int{1001, 1002},
			serverDatasetIdsOrder: []int{1002, 1001},
			format:                datastream.FormatTypeStructured,
			expectNonEmptyPlan:    true,
		},
		"id change in structured config": {
			preConfig:             "testdata/TestResourceStream/dataset_ids_diff/structured_config.tf",
			fileDatasetIdsOrder:   []int{1001, 1002},
			serverDatasetIdsOrder: []int{1002, 1003},
			format:                datastream.FormatTypeStructured,
			expectNonEmptyPlan:    true,
		},
	}

	for name, test := range tests {

		streamConfiguration := datastream.StreamConfiguration{
			ActivateNow: false,
			Config: datastream.Config{
				Format: test.format,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: DefaultUploadFilePrefix,
				UploadFileSuffix: DefaultUploadFileSuffix,
			},
			Connectors: []datastream.AbstractConnector{
				&datastream.SplunkConnector{
					CompressLogs:        false,
					ConnectorName:       "splunk_test_connector_name",
					EventCollectorToken: "splunk_event_collector_token",
					URL:                 "splunk_url",
				},
			},
			ContractID:      "test_contract",
			DatasetFieldIDs: test.fileDatasetIdsOrder,
			GroupID:         tools.IntPtr(1337),
			PropertyIDs:     []int{1},
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}

		if test.format == datastream.FormatTypeStructured {
			streamConfiguration.Config.Delimiter = datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace)
		}

		createStreamRequest := datastream.CreateStreamRequest{
			StreamConfiguration: streamConfiguration,
		}

		createStreamResponse := &datastream.StreamUpdate{
			StreamVersionKey: datastream.StreamVersionKey{
				StreamID:        streamID,
				StreamVersionID: 1,
			},
		}

		getStreamRequest := datastream.GetStreamRequest{
			StreamID: streamID,
		}

		getStreamResponse := &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				{
					ConnectorType: datastream.ConnectorTypeSplunk,
					CompressLogs:  false,
					ConnectorName: "splunk_test_connector_name",
					URL:           "splunk_url",
				},
			},
			ContractID: streamConfiguration.ContractID,
			Datasets: []datastream.DataSets{
				{
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID: test.serverDatasetIdsOrder[0],
							Order:          0,
						},
						{
							DatasetFieldID: test.serverDatasetIdsOrder[1],
							Order:          1,
						},
					},
				},
			},
			GroupID: *streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: 2,
			TemplateName:    streamConfiguration.TemplateName,
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
				Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString(test.preConfig),
							ExpectNonEmptyPlan: test.expectNonEmptyPlan,
							Check: resource.ComposeTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields_ids.#", "2"),
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields_ids.0", strconv.Itoa(test.serverDatasetIdsOrder[0])),
								resource.TestCheckResourceAttr("akamai_datastream.splunk_stream", "dataset_fields_ids.1", strconv.Itoa(test.serverDatasetIdsOrder[1])),
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
		ActivateNow: false,
		Config: datastream.Config{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				TimeInSec: datastream.TimeInSec30,
			},
			UploadFilePrefix: DefaultUploadFilePrefix,
			UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID:      "test_contract",
		DatasetFieldIDs: []int{1001},
		GroupID:         tools.IntPtr(1337),
		PropertyIDs:     []int{1},
		StreamName:      "test_stream",
		StreamType:      datastream.StreamTypeRawLogs,
		TemplateName:    datastream.TemplateNameEdgeLogs,
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Connectors = []datastream.AbstractConnector{
			connector,
		}
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
		}
	}

	responseFactory := func(connector datastream.ConnectorDetails) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				connector,
			},
			ContractID: streamConfiguration.ContractID,
			Datasets: []datastream.DataSets{
				{
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID: 1001,
							Order:          0,
						},
					},
				},
			},
			GroupID: *streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: 2,
			TemplateName:    streamConfiguration.TemplateName,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 1,
		},
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.ConnectorDetails
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"splunk": {
			Filename: "custom_headers_splunk.tf",
			Connector: &datastream.SplunkConnector{
				CompressLogs:        false,
				ConnectorName:       "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				URL:                 "splunk_url",
				CustomHeaderName:    "custom_header_name",
				CustomHeaderValue:   "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:     datastream.ConnectorTypeSplunk,
				CompressLogs:      false,
				ConnectorName:     "splunk_test_connector_name",
				URL:               "splunk_url",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.connector_name", "splunk_test_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.event_collector_token", "splunk_event_collector_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.url", "splunk_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.custom_header_value", "custom_header_value"),
			},
		},
		"https": {
			Filename: "custom_headers_https.tf",
			Connector: &datastream.CustomHTTPSConnector{
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				ConnectorName:      "HTTPS connector name",
				Password:           "password",
				URL:                "https_connector_url",
				UserName:           "username",
				ContentType:        "content_type",
				CustomHeaderName:   "custom_header_name",
				CustomHeaderValue:  "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:      datastream.ConnectorTypeHTTPS,
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				ConnectorName:      "HTTPS connector name",
				URL:                "https_connector_url",
				ContentType:        "content_type",
				CustomHeaderName:   "custom_header_name",
				CustomHeaderValue:  "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.authentication_type", "BASIC"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.connector_name", "HTTPS connector name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.compress_logs", "true"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.custom_header_name", "custom_header_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.custom_header_value", "custom_header_value"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.url", "https_connector_url"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.user_name", "username"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.password", "password"),
			},
		},
		"sumologic": {
			Filename: "custom_headers_sumologic.tf",
			Connector: &datastream.SumoLogicConnector{
				CollectorCode:     "collector_code",
				CompressLogs:      true,
				ConnectorName:     "Sumologic connector name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:     datastream.ConnectorTypeSumoLogic,
				CompressLogs:      true,
				ConnectorName:     "Sumologic connector name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "Sumologic connector name"),
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
				ConnectorName:     "loggly_connector_name",
				Endpoint:          "endpoint",
				AuthToken:         "auth_token",
				Tags:              "tag1,tag2,tag3",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:     datastream.ConnectorTypeLoggly,
				ConnectorName:     "loggly_connector_name",
				Endpoint:          "endpoint",
				Tags:              "tag1,tag2,tag3",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "loggly_connector.0.connector_name", "loggly_connector_name"),
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
				ConnectorName:     "new_relic_connector_name",
				Endpoint:          "endpoint",
				AuthToken:         "auth_token",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:     datastream.ConnectorTypeNewRelic,
				ConnectorName:     "new_relic_connector_name",
				Endpoint:          "endpoint",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "new_relic_connector.0.connector_name", "new_relic_connector_name"),
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
				ConnectorName:     "elasticsearch_connector_name",
				Endpoint:          "endpoint",
				IndexName:         "index_name",
				UserName:          "user_name",
				Password:          "password",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:     datastream.ConnectorTypeElasticsearch,
				ConnectorName:     "elasticsearch_connector_name",
				Endpoint:          "endpoint",
				IndexName:         "index_name",
				ContentType:       "content_type",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.connector_name", "elasticsearch_connector_name"),
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
			}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/custom_headers/%s", test.Filename)),
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
		ActivateNow: false,
		Config: datastream.Config{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				TimeInSec: datastream.TimeInSec30,
			},
			UploadFilePrefix: DefaultUploadFilePrefix,
			UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID:      "test_contract",
		DatasetFieldIDs: []int{1001},
		GroupID:         tools.IntPtr(1337),
		PropertyIDs:     []int{1},
		StreamName:      "test_stream",
		StreamType:      datastream.StreamTypeRawLogs,
		TemplateName:    datastream.TemplateNameEdgeLogs,
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Connectors = []datastream.AbstractConnector{
			connector,
		}
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
		}
	}

	responseFactory := func(connector datastream.ConnectorDetails) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				connector,
			},
			ContractID: streamConfiguration.ContractID,
			Datasets: []datastream.DataSets{
				{
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID: 1001,
							Order:          0,
						},
					},
				},
			},
			GroupID: *streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: 2,
			TemplateName:    streamConfiguration.TemplateName,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 1,
		},
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.ConnectorDetails
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"splunk_mtls": {
			Filename: "mtls_splunk.tf",
			Connector: &datastream.SplunkConnector{
				CompressLogs:        false,
				ConnectorName:       "splunk_test_connector_name",
				EventCollectorToken: "splunk_event_collector_token",
				URL:                 "splunk_url",
				TLSHostname:         "tls_hostname",
				CACert:              "ca_cert",
				ClientCert:          "client_cert",
				ClientKey:           "client_key",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType: datastream.ConnectorTypeSplunk,
				CompressLogs:  false,
				ConnectorName: "splunk_test_connector_name",
				URL:           "splunk_url",
				TLSHostname:   "tls_hostname",
				MTLS:          "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.connector_name", "splunk_test_connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.event_collector_token", "splunk_event_collector_token"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "splunk_connector.0.url", "splunk_url"),
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
				ConnectorName:      "HTTPS connector name",
				Password:           "password",
				URL:                "https_connector_url",
				UserName:           "username",
				ContentType:        "content_type",
				TLSHostname:        "tls_hostname",
				CACert:             "ca_cert",
				ClientCert:         "client_cert",
				ClientKey:          "client_key",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:      datastream.ConnectorTypeHTTPS,
				AuthenticationType: datastream.AuthenticationTypeBasic,
				CompressLogs:       true,
				ConnectorName:      "HTTPS connector name",
				URL:                "https_connector_url",
				ContentType:        "content_type",
				TLSHostname:        "tls_hostname",
				MTLS:               "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.authentication_type", "BASIC"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.connector_name", "HTTPS connector name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.compress_logs", "true"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.content_type", "content_type"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "https_connector.0.url", "https_connector_url"),
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
				ConnectorName: "elasticsearch_connector_name",
				Endpoint:      "endpoint",
				IndexName:     "index_name",
				UserName:      "user_name",
				Password:      "password",
				ContentType:   "content_type",
				TLSHostname:   "tls_hostname",
				CACert:        "ca_cert",
				ClientCert:    "client_cert",
				ClientKey:     "client_key",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType: datastream.ConnectorTypeElasticsearch,
				CompressLogs:  true,
				ConnectorName: "elasticsearch_connector_name",
				Endpoint:      "endpoint",
				IndexName:     "index_name",
				ContentType:   "content_type",
				TLSHostname:   "tls_hostname",
				MTLS:          "Enabled",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "elasticsearch_connector.0.connector_name", "elasticsearch_connector_name"),
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
			}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/mtls/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}
}

func TestUrlSuppressor(t *testing.T) {

	streamConfigurationFactory := func(connector datastream.AbstractConnector) datastream.StreamConfiguration {
		return datastream.StreamConfiguration{
			ActivateNow: false,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "ak",
				UploadFileSuffix: "ds",
			},
			Connectors: []datastream.AbstractConnector{
				connector,
			},
			ContractID:      "test_contract",
			DatasetFieldIDs: []int{1001},
			GroupID:         tools.IntPtr(1337),
			PropertyIDs:     []int{1},
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationFactory(connector),
		}
	}

	updateStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.UpdateStreamRequest {
		req := datastream.UpdateStreamRequest{
			StreamID:            streamID,
			StreamConfiguration: streamConfigurationFactory(connector),
		}
		req.StreamConfiguration.GroupID = nil
		return req
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 1,
		},
	}

	responseFactory := func(connector datastream.ConnectorDetails) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
			},
			Connectors: []datastream.ConnectorDetails{
				connector,
			},
			ContractID: "test_contract",
			Datasets: []datastream.DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          1001,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   0,
						},
					},
				},
			},
			GroupID: 1337,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			StreamVersionID: 1,
			TemplateName:    datastream.TemplateNameEdgeLogs,
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
					ConnectorName: "connector_name",
					Endpoint:      "endpoint/?/?",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeSumoLogic,
						CompressLogs:  true,
						ConnectorName: "connector_name",
						Endpoint:      "endpoint", //api returns stripped url
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/idempotency/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config:   loadFixtureString("testdata/TestResourceStream/urlSuppressor/idempotency/create_stream.tf"),
					PlanOnly: true,
				},
			},
		},
		"update endpoint field": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					ConnectorName: "connector_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeSumoLogic,
						CompressLogs:  true,
						ConnectorName: "connector_name",
						Endpoint:      "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					ConnectorName: "connector_name",
					Endpoint:      "endpoint_updated",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeSumoLogic,
						CompressLogs:  true,
						ConnectorName: "connector_name",
						Endpoint:      "endpoint_updated",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/update_endpoint_field/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/update_endpoint_field/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
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
					ConnectorName: "connector_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeSumoLogic,
						CompressLogs:  true,
						ConnectorName: "connector_name",
						Endpoint:      "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode:     "collector_code",
					CompressLogs:      true,
					ConnectorName:     "connector_name",
					Endpoint:          "endpoint",
					CustomHeaderName:  "custom_header_name",
					CustomHeaderValue: "custom_header_value",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType:     datastream.ConnectorTypeSumoLogic,
						CompressLogs:      true,
						ConnectorName:     "connector_name",
						Endpoint:          "endpoint",
						CustomHeaderName:  "custom_header_name",
						CustomHeaderValue: "custom_header_value",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/adding_fields/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/adding_fields/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_name", "custom_header_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.custom_header_value", "custom_header_value"),
					),
				},
				{
					Config:   loadFixtureString("testdata/TestResourceStream/urlSuppressor/adding_fields/update_stream.tf"),
					PlanOnly: true,
				},
			},
		},
		"change connector": {
			Init: func(m *datastream.Mock) {
				m.On("CreateStream", mock.Anything, createStreamRequestFactory(&datastream.SumoLogicConnector{
					CollectorCode: "collector_code",
					CompressLogs:  true,
					ConnectorName: "connector_name",
					Endpoint:      "endpoint",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeSumoLogic,
						CompressLogs:  true,
						ConnectorName: "connector_name",
						Endpoint:      "endpoint",
					}), nil).Times(3)

				m.On("UpdateStream", mock.Anything, updateStreamRequestFactory(&datastream.DatadogConnector{
					AuthToken:     "auth_token",
					ConnectorName: "connector_name",
					URL:           "url",
				})).Return(updateStreamResponse, nil)

				m.On("GetStream", mock.Anything, mock.Anything).
					Return(responseFactory(datastream.ConnectorDetails{
						ConnectorType: datastream.ConnectorTypeDataDog,
						ConnectorName: "connector_name",
						URL:           "url",
					}), nil)
			},
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/change_connector/create_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.collector_code", "collector_code"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "sumologic_connector.0.endpoint", "endpoint"),
					),
				},
				{
					Config: loadFixtureString("testdata/TestResourceStream/urlSuppressor/change_connector/update_stream.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.connector_name", "connector_name"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.auth_token", "auth_token"),
						resource.TestCheckResourceAttr("akamai_datastream.s", "datadog_connector.0.url", "url"),
					),
				},
				{
					Config:   loadFixtureString("testdata/TestResourceStream/urlSuppressor/change_connector/update_stream.tf"),
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
			}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(m, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps:     test.Steps,
				})

				m.AssertExpectations(t)
			})
		})
	}
}

func TestConnectors(t *testing.T) {
	streamConfiguration := datastream.StreamConfiguration{
		ActivateNow: false,
		Config: datastream.Config{
			Format: datastream.FormatTypeJson,
			Frequency: datastream.Frequency{
				TimeInSec: datastream.TimeInSec30,
			},
			UploadFilePrefix: DefaultUploadFilePrefix,
			UploadFileSuffix: DefaultUploadFileSuffix,
		},
		ContractID:      "test_contract",
		DatasetFieldIDs: []int{1001},
		GroupID:         tools.IntPtr(1337),
		PropertyIDs:     []int{1},
		StreamName:      "test_stream",
		StreamType:      datastream.StreamTypeRawLogs,
		TemplateName:    datastream.TemplateNameEdgeLogs,
	}

	createStreamRequestFactory := func(connector datastream.AbstractConnector) datastream.CreateStreamRequest {
		streamConfigurationWithConnector := streamConfiguration
		streamConfigurationWithConnector.Connectors = []datastream.AbstractConnector{
			connector,
		}
		return datastream.CreateStreamRequest{
			StreamConfiguration: streamConfigurationWithConnector,
		}
	}

	responseFactory := func(connector datastream.ConnectorDetails) *datastream.DetailedStreamVersion {
		return &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusInactive,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				connector,
			},
			ContractID: streamConfiguration.ContractID,
			Datasets: []datastream.DataSets{
				{
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID: 1001,
							Order:          0,
						},
					},
				},
			},
			GroupID: *streamConfiguration.GroupID,
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
			},
			StreamID:        streamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: 2,
			TemplateName:    streamConfiguration.TemplateName,
		}
	}

	getStreamRequest := datastream.GetStreamRequest{
		StreamID: streamID,
	}

	updateStreamResponse := &datastream.StreamUpdate{
		StreamVersionKey: datastream.StreamVersionKey{
			StreamID:        streamID,
			StreamVersionID: 1,
		},
	}

	tests := map[string]struct {
		Filename   string
		Response   datastream.ConnectorDetails
		Connector  datastream.AbstractConnector
		TestChecks []resource.TestCheckFunc
	}{
		"azure": {
			Filename: "azure.tf",
			Connector: &datastream.AzureConnector{
				AccessKey:     "access_key",
				AccountName:   "account_name",
				ConnectorName: "connector_name",
				ContainerName: "container_name",
				Path:          "path",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType: datastream.ConnectorTypeAzure,
				AccountName:   "account_name",
				ConnectorName: "connector_name",
				ContainerName: "container_name",
				Path:          "path",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.account_name", "account_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.connector_name", "connector_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.container_name", "container_name"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "azure_connector.0.path", "path"),
			},
		},
		"gcs": {
			Filename: "gcs.tf",
			Connector: &datastream.GCSConnector{
				Bucket:             "bucket",
				ConnectorName:      "connector_name",
				Path:               "path",
				PrivateKey:         "private_key",
				ProjectID:          "project_id",
				ServiceAccountName: "service_account_name",
			},
			Response: datastream.ConnectorDetails{
				ConnectorType:      datastream.ConnectorTypeGcs,
				Bucket:             "bucket",
				ConnectorName:      "connector_name",
				Path:               "path",
				ProjectID:          "project_id",
				ServiceAccountName: "service_account_name",
			},
			TestChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.#", "1"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.bucket", "bucket"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.compress_logs", "false"),
				resource.TestCheckResourceAttr("akamai_datastream.s", "gcs_connector.0.connector_name", "connector_name"),
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

			client.On("DeleteStream", mock.Anything, mock.Anything).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourceStream/connectors/%s", test.Filename)),
							Check:  resource.ComposeTestCheckFunc(test.TestChecks...),
						},
					},
				})

				client.AssertExpectations(t)
			})
		})
	}
}
